// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/crypto"
	"github.com/clearmatics/autonity/consensus/tendermint/events"
	"github.com/clearmatics/autonity/contracts/autonity"
	"github.com/clearmatics/autonity/core/types"
	autonitycrypto "github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/rlp"
)

// Start implements core.Tendermint.Start
func (c *core) Start(ctx context.Context, contract *autonity.Contract) {
	// Set the autonity contract
	c.autonityContract = contract
	ctx, c.cancel = context.WithCancel(ctx)

	c.subscribeEvents()

	// core.height needs to be set beforehand for unmined block's logic.
	lastBlockMined, _ := c.backend.LastCommittedProposal()
	c.setHeight(new(big.Int).Add(lastBlockMined.Number(), common.Big1))
	// We need a separate go routine to keep c.latestPendingUnminedBlock up to date
	go c.handleNewUnminedBlockEvent(ctx)

	// Tendermint Finite State Machine discrete event loop
	go c.mainEventLoop(ctx)

	go c.backend.HandleUnhandledMsgs(ctx)
}

// Stop implements core.Engine.Stop
func (c *core) Stop() {
	c.logger.Info("stopping tendermint.core", "addr", c.address.String())

	_ = c.proposeTimeout.stopTimer()
	_ = c.prevoteTimeout.stopTimer()
	_ = c.precommitTimeout.stopTimer()

	c.cancel()

	c.stopFutureProposalTimer()
	c.unsubscribeEvents()

	// Ensure all event handling go routines exit
	<-c.stopped
	<-c.stopped
	<-c.stopped
}

func (c *core) subscribeEvents() {
	s := c.backend.Subscribe(events.MessageEvent{}, backlogEvent{})
	c.messageEventSub = s

	s1 := c.backend.Subscribe(events.NewUnminedBlockEvent{})
	c.newUnminedBlockEventSub = s1

	s2 := c.backend.Subscribe(TimeoutEvent{})
	c.timeoutEventSub = s2

	s3 := c.backend.Subscribe(events.CommitEvent{})
	c.committedSub = s3

	s4 := c.backend.Subscribe(events.SyncEvent{})
	c.syncEventSub = s4
}

// Unsubscribe all messageEventSub
func (c *core) unsubscribeEvents() {
	c.messageEventSub.Unsubscribe()
	c.newUnminedBlockEventSub.Unsubscribe()
	c.timeoutEventSub.Unsubscribe()
	c.committedSub.Unsubscribe()
	c.syncEventSub.Unsubscribe()
}

// TODO: update all of the TypeMuxSilent to event.Feed and should not use backend.EventMux for core internal messageEventSub: backlogEvent, TimeoutEvent

func (c *core) handleNewUnminedBlockEvent(ctx context.Context) {
eventLoop:
	for {
		select {
		case e, ok := <-c.newUnminedBlockEventSub.Chan():
			if !ok {
				break eventLoop
			}
			newUnminedBlockEvent := e.Data.(events.NewUnminedBlockEvent)
			pb := &newUnminedBlockEvent.NewUnminedBlock
			c.storeUnminedBlockMsg(pb)
		case <-ctx.Done():
			c.logger.Info("handleNewUnminedBlockEvent is stopped", "event", ctx.Err())
			break eventLoop
		}
	}

	c.stopped <- struct{}{}
}

func (c *core) mainEventLoop(ctx context.Context) {
	// Start a new round from last height + 1
	c.startRound(ctx, 0)

	go c.syncLoop(ctx)

eventLoop:
	for {
		select {
		case ev, ok := <-c.messageEventSub.Chan():
			if !ok {
				break eventLoop
			}
			// A real ev arrived, process interesting content
			switch e := ev.Data.(type) {
			case events.MessageEvent:
				if len(e.Payload) == 0 {
					c.logger.Error("core.mainEventLoop Get message(MessageEvent) empty payload")
				}

				if err := c.handleMsg(ctx, e.Payload); err != nil {
					c.logger.Debug("core.mainEventLoop Get message(MessageEvent) payload failed", "err", err)
					continue
				}
				c.backend.Gossip(ctx, c.committeeSet().Committee(), e.Payload)
			case proposalEvent:
				err := c.handleProposal(ctx, e.proposal)
				if err != nil {
					c.logger.Debug("core.mainEventLoop handleProposal message failed", "err", err)
				}
			}
		case ev, ok := <-c.timeoutEventSub.Chan():
			if !ok {
				break eventLoop
			}
			if timeoutE, ok := ev.Data.(TimeoutEvent); ok {
				switch timeoutE.step {
				case msgProposal:
					c.handleTimeoutPropose(ctx, timeoutE)
				case msgPrevote:
					c.handleTimeoutPrevote(ctx, timeoutE)
				case msgPrecommit:
					c.handleTimeoutPrecommit(ctx, timeoutE)
				}
			}
		case ev, ok := <-c.committedSub.Chan():
			if !ok {
				break eventLoop
			}
			switch ev.Data.(type) {
			case events.CommitEvent:
				c.handleCommit(ctx)
			}
		case <-ctx.Done():
			c.logger.Info("mainEventLoop is stopped", "event", ctx.Err())
			break eventLoop
		}
	}

	c.stopped <- struct{}{}
}

func (c *core) syncLoop(ctx context.Context) {
	/*
		this method is responsible for asking the network to send us the current consensus state
		and to process sync queries events.
	*/
	timer := time.NewTimer(10 * time.Second)

	round := c.Round()
	height := c.Height()

	// Ask for sync when the engine starts
	c.backend.AskSync(c.lastHeader)

eventLoop:
	for {
		select {
		case <-timer.C:
			currentRound := c.Round()
			currentHeight := c.Height()

			// we only ask for sync if the current view stayed the same for the past 10 seconds
			if currentHeight.Cmp(height) == 0 && currentRound == round {
				c.backend.AskSync(c.lastHeader)
			}
			round = currentRound
			height = currentHeight
			timer = time.NewTimer(10 * time.Second)

		case ev, ok := <-c.syncEventSub.Chan():
			if !ok {
				break eventLoop
			}
			event := ev.Data.(events.SyncEvent)
			c.logger.Info("Processing sync message", "from", event.Addr)
			c.backend.SyncPeer(event.Addr)
		case <-ctx.Done():
			c.logger.Info("syncLoop is stopped", "event", ctx.Err())
			break eventLoop
		}
	}

	c.stopped <- struct{}{}
}

// sendEvent sends event to mux
func (c *core) sendEvent(ev interface{}) {
	c.backend.Post(ev)
}

type consensusMessage struct {
	msgType    uint8
	height     uint64
	round      int64
	value      common.Hash
	validRound int64
}

func (c *core) handleMsg(ctx context.Context, payload []byte) error {

	/*
		Basic validity checks
	*/

	m := new(Message)

	// Set the hash on the message so that it can be used for indexing.
	m.Hash = common.BytesToHash(autonitycrypto.Keccak256(payload))

	// Check we haven't already processed this message
	if c.msgCache.Message(m.Hash) != nil {
		// Message was already processed
		return nil
	}

	// Decode message
	err := rlp.DecodeBytes(payload, m)
	if err != nil {
		return err
	}

	var proposal Proposal
	var preVote Vote
	var preCommit Vote
	var conMsg *consensusMessage
	switch m.Code {
	case msgProposal:
		err := m.Decode(&proposal)
		if err != nil {
			return errFailedDecodeProposal
		}

		valueHash := proposal.ProposalBlock.Hash()
		conMsg = &consensusMessage{
			msgType:    uint8(m.Code),
			height:     proposal.Height.Uint64(),
			round:      proposal.Round,
			value:      valueHash,
			validRound: proposal.ValidRound,
		}

		err = c.msgCache.addMessage(m, conMsg)
		if err != nil {
			// could be multipe proposal messages from the same proposer
			return err
		}
		c.msgCache.addValue(valueHash, proposal.ProposalBlock)

	case msgPrevote:
		err := m.Decode(&preVote)
		if err != nil {
			return errFailedDecodePrevote
		}
		conMsg = &consensusMessage{
			msgType: uint8(m.Code),
			height:  preVote.Height.Uint64(),
			round:   preVote.Round,
			value:   preVote.ProposedBlockHash,
		}

		err = c.msgCache.addMessage(m, conMsg)
		if err != nil {
			// could be multipe precommits from same validator
			return err
		}
	case msgPrecommit:
		err := m.Decode(&preCommit)
		if err != nil {
			return errFailedDecodePrecommit
		}
		// Check the committed seal matches the block hash if its a precommit.
		// If not we ignore the message.
		//
		// Note this method does not make use of any blockchain state, so it is
		// safe to call it now. In fact it only uses the logger of c so I think
		// it could easily be detached from c.
		err = c.verifyCommittedSeal(m.Address, append([]byte(nil), m.CommittedSeal...), preCommit.ProposedBlockHash, preCommit.Round, preCommit.Height)
		if err != nil {
			return err
		}
		conMsg = &consensusMessage{
			msgType: uint8(m.Code),
			height:  preCommit.Height.Uint64(),
			round:   preCommit.Round,
			value:   preCommit.ProposedBlockHash,
		}

		err = c.msgCache.addMessage(m, conMsg)
		if err != nil {
			// could be multipe precommits from same validator
			return err
		}
	default:
		return fmt.Errorf("unrecognised consensus message code %q", m.Code)
	}

	// If this message is for a future height then we cannot validate it
	// because we lack the relevant header, we will process it when we reach
	// that height. If it is for a previous height then we are not intersted in
	// it. But it has been added to the msg cache in case other peers would
	// like to sync it.
	if conMsg.height != c.Height().Uint64() {
		// Nothing to do here
		return nil
	}

	return handleCurrentHeightMessage(m, conMsg)

}

func (c *core) handleCurrentHeightMessage(m *Message, cm *consensusMessage) error {
	/*
		Domain specific validity checks, now we know that we are at the same
		height as this message we can rely on lastHeader.
	*/

	// Check that the message came from a committee member, if not we ignore it.
	if c.lastHeader.CommitteeMember(m.Address) == nil {
		// TODO turn this into an error type that can be checked for at a
		// higher level to close the connection to this peer.
		return fmt.Errorf("received message from non committee member: %v", m)
	}

	payload, err := m.PayloadNoSig()
	if err != nil {
		return err
	}

	// Again we ignore messges with invalid signatures, they cannot be trusted.
	// TODO make crypto.CheckValidatorSignature accept Message so that it can
	// handle generating the payload and checking it with the sig and address.
	address, err := crypto.CheckValidatorSignature(c.lastHeader, payload, m.Signature)
	if err != nil {
		return err
	}

	if address != m.Address {
		// TODO why is Address even a field of Message when the address can be derived?
		return fmt.Errorf("address in message %q and address derived from signature %q don't match", m.Address, address)
	}

	switch m.Code {
	case msgProposal:
		// We ignore proposals from non proposers
		if !c.isProposerMsg(proposal.Round, m.Address) {
			c.logger.Warn("Ignore proposal messages from non-proposer")
			return errNotFromProposer

			// TODO verify proposal here.
			//
			// If we are introducing time into the mix then what we are saying
			// is that we don't expect different participants' clocks to drift
			// out of sync more than some delta. And if they do then we don't
			// expect consensus to work.
			//
			// So in the case that clocks drift too far out of sync and say a
			// node considers a proposal invalid that 2f+1 other nodes
			// precommit for that node becomes stuck and can only continue in
			// consensus by re-syncing the blocks.
			//
			// So in verifying the proposal wrt time we should verify once
			// within reasonable clock sync bounds and then set the validity
			// based on that and never re-process the message again.
		}
	}

	c.msgCache.setValid(m.Hash)
	c.checkUponConditions(cm)

	return nil
}

var (
	voteForNil   bool        = true
	voteForValue bool        = false
	nilValue     common.Hash = common.Hash{}
)

func (c *core) checkUponConditions(cm *consensusMessage) {
	r := c.Round()
	h := c.Height()
	lh := c.lastHeader
	s := c.step

	// Some of the checks in these upon conditions are omitted because they have alrady been checked.
	//
	// - We do not check height because we only execute this code when the
	// message height matches the current height.
	//
	// - We do not check whether the message comes from a proposer since this
	// is checkded before calling this method and we do not process proposals
	// from non proposers.
	//
	// - We do not whether proposal values are valid since this is checked
	// before calling this method and we do not process proposals that are not
	// valid.

	// Line 22
	if cm.msgType == uint8(msgProposal) && cm.round == r && cm.validRound == -1 && c.step == propose {
		if c.lockedRound == -1 || c.lockedValue.Hash() == cm.value {
			c.sendPrevote(nil, voteForValue)
		} else {
			c.sendPrevote(nil, voteForNil)
		}
	}

	// Line 28
	if cm.msgType == uint8(msgProposal) && cm.round == r && c.msgCache.prevoteQuorum(&cm.value, cm.validRound, lh) && s == propose && (cm.validRound >= 0 && cm.validRound < r) {
		if c.lockedRound <= cm.validRound || c.lockedValue.Hash() == cm.value {
			c.sendPrevote(nil, voteForValue)
		} else {
			c.sendPrevote(nil, voteForNil)
		}
	}

	// Line 34
	if c.msgCache.prevoteQuorum(nil, cm.round, lh) && s == prevote && !c.line34Executed {
		c.prevoteTimeout.scheduleTimeout(c.timeoutPrevote(r), r, h, c.onTimeoutPrecommit)
	}

	// Line 36
	if cm.msgType == uint8(msgProposal) && cm.round == r && c.msgCache.prevoteQuorum(&cm.value, r, lh) && s >= prevote && !c.line36Executed {
		block := c.msgCache.value(cm.value) // TODO remove references to block from core
		if s == prevote {
			c.lockedValue = block
			c.lockedRound = r
			c.sendPrecommit(nil, voteForValue)
			s = precommit
			c.step = s
		}
		c.validValue = block
		c.validRound = cm.round
	}

	// Line 44
	if c.msgCache.prevoteQuorum(&nilValue, r, lh) && s == prevote {
		c.sendPrecommit(nil, voteForValue)
		s = precommit
		c.step = s
	}

	// Line 47
	if c.msgCache.precommitQuorum(nil, cm.round, lh) && !c.line47Executed {
		c.precommitTimeout.scheduleTimeout(c.timeoutPrecommit(r), r, h, c.onTimeoutPrecommit)
	}

	// Line 49
	if cm.msgType == uint8(msgProposal) && c.msgCache.precommitQuorum(&cm.value, cm.round, lh) {
		block := c.msgCache.value(cm.value) // TODO remove references to block from core
		if s == prevote {
			c.commit(block)
			c.setHeight(block.NumberU64() + 1)
			c.lockedRound = -1
			c.lockedValue = nil
			c.validRound = -1
			c.validValue = nil

			// Not quite sure how to start the round nicely
			// need to ensure that we don't stack overflow in the case that the
			// next height messages are sufficient for consensus when we
			// process them and so on and so on.  So I need to set the start
			// round states and then queue the messages for processing. And I
			// need to ensure that I get a list of messages to process in an
			// atomic step from the msg cache so that I don't end up trying to
			// process the same message twice.
		}
	}

	// TODO need to re-write all these rules to take into account all message
	// types that could match eg line 49 could be a propsal or a precommit.
}

func (c *core) handleFutureRoundMsg(ctx context.Context, msg *Message, sender types.CommitteeMember) {
	// Decoding functions can't fail here
	msgRound, err := msg.Round()
	if err != nil {
		c.logger.Error("handleFutureRoundMsg msgRound", "err", err)
		return
	}
	if _, ok := c.futureRoundChange[msgRound]; !ok {
		c.futureRoundChange[msgRound] = make(map[common.Address]uint64)
	}
	c.futureRoundChange[msgRound][sender.Address] = sender.VotingPower.Uint64()

	var totalFutureRoundMessagesPower uint64
	for _, power := range c.futureRoundChange[msgRound] {
		totalFutureRoundMessagesPower += power
	}

	if totalFutureRoundMessagesPower > c.committeeSet().F() {
		c.logger.Info("Received ceil(N/3) - 1 messages power for higher round", "New round", msgRound)
		c.startRound(ctx, msgRound)
	}
}
