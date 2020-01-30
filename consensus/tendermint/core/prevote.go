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
	"math/big"

	"github.com/clearmatics/autonity/common"
)

func (c *core) sendPrevote(ctx context.Context, height *big.Int, round *big.Int, hash common.Hash) {
	logger := c.logger.New("step", c.getStep())

	var prevote = Vote{
		Height:            height,
		Round:             round,
		ProposedBlockHash: hash,
	}

	encodedVote, err := Encode(&prevote)
	if err != nil {
		logger.Error("Failed to encode", "subject", prevote)
		return
	}

	c.logPrevoteMessageEvent("MessageEvent(Prevote): Sent", prevote, c.address.String(), "broadcast")

	c.sentPrevote = true
	c.broadcast(ctx, &Message{
		Code:          msgPrevote,
		Msg:           encodedVote,
		Address:       c.address,
		CommittedSeal: []byte{},
	})
}

func (c *core) handlePrevote(ctx context.Context, msg *Message) error {
	var preVote Vote
	err := msg.Decode(&preVote)
	if err != nil {
		return errFailedDecodePrevote
	}

	// Check for nil values
	if preVote.Height == nil || preVote.Round == nil {
		return errInvalidMessage
	}

	// Ensure prevote is for current height
	if err = c.checkMessage(preVote.Round, preVote.Height, prevote); err != nil {
		return err
	}

	// If we already have the prevote do nothing
	if c.hasVote(preVote, msg) {
		return nil
	}

	curR := c.getRound().Int64()
	curH := c.getHeight().Int64()
	prevoteHash := preVote.ProposedBlockHash

	// The prevote doesn't exists in our current round state, so add it
	if prevotes := c.getPrevotesSet(preVote.Round.Int64()); prevotes == nil {
		c.setPrevotesSet(preVote.Round.Int64())
	}
	prevotes := c.getPrevotesSet(preVote.Round.Int64())
	prevotes.Add(prevoteHash, *msg)

	c.logPrevoteMessageEvent("MessageEvent(Prevote): Received", preVote, msg.Address.String(), c.address.String())

	roundCmp := preVote.Round.Cmp(c.getRound())
	if roundCmp < 0 {
		return c.checkForOldProposal(ctx, curR)
	} else if roundCmp > 0 {
		// TODO: check if validator needs to move to a future round
	} else {
		// preVote.Round.Int64()==curR
		c.checkForPrevoteTimeout(curR, curH)
		if err := c.checkForQuorumPrevotes(ctx, curR); err != nil {
			return err
		}
		if err := c.checkForQuorumPrevotesNil(ctx, curR); err != nil {
			return err
		}
	}
	return nil
}

func (c *core) logPrevoteMessageEvent(message string, prevote Vote, from, to string) {
	currentRound := c.getRound().Int64()
	currentProposalHash := common.Hash{}

	proposalMS := c.getProposalSet(currentRound)
	if proposalMS != nil {
		currentProposalHash = proposalMS.proposal().ProposalBlock.Hash()
	}

	prevotes := c.getPrevotesSet(currentRound)
	if prevotes == nil {
		return
	}
	c.logger.Debug(message,
		"from", from,
		"to", to,
		"currentHeight", c.getHeight(),
		"msgHeight", prevote.Height,
		"currentRound", c.getRound(),
		"msgRound", prevote.Round,
		"currentStep", c.getStep(),
		"isProposer", c.isProposerForR(c.getRound().Int64(), c.address),
		"currentProposer", c.valSet.GetProposer(),
		"isNilMsg", prevote.ProposedBlockHash == common.Hash{},
		"hash", prevote.ProposedBlockHash,
		"type", "Prevote",
		"totalVotes", prevotes.TotalSize(),
		"totalNilVotes", prevotes.NilVotesSize(),
		"quorumReject", c.quorum(prevotes.NilVotesSize()),
		"totalNonNilVotes", prevotes.VotesSize(currentProposalHash),
		"quorumAccept", c.quorum(prevotes.VotesSize(currentProposalHash)),
	)
}