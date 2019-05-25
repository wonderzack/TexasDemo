package texas_holdem

import (
	"fmt"
	"sort"
)

type IBetStatus interface {
	SeatID() int32
	WinVal() uint32   // GiveUp 的话winval 为 0 否则等于 (h *Hand) FinalLevel() + 1
	BetAmount() int64 // 下注金额
}
type SeatID2WinAmount map[int32]int64

func DistributePond(inputs []IBetStatus) SeatID2WinAmount {
	seatID2Win := make(SeatID2WinAmount)

	betStatusSlice := make(BetStatusSlice, len(inputs))
	for index, input := range inputs {
		betStatusSlice[index] = &BetStatus{
			seatID:    input.SeatID(),
			winVal:    input.WinVal(),
			betAmount: input.BetAmount(),
		}
	}

	sort.Sort(betStatusSlice)

	winners := make([]int32, 0, len(betStatusSlice))
	for {
		// 删除 BetAmount=0 的项 并找出最小的 BetAmount
		cursor := 0
		for i := 0; i < len(betStatusSlice); i++ {
			eachStatus := betStatusSlice[i]
			if eachStatus.betAmount != 0 {
				betStatusSlice[cursor] = eachStatus
				cursor++
			}
		}
		betStatusSlice = betStatusSlice[:cursor]

		if len(betStatusSlice) == 0 {
			break
		}

		if len(betStatusSlice) == 1 {
			seatID2Win[betStatusSlice[0].seatID] += betStatusSlice[0].betAmount
			break
		}

		firstWinner := betStatusSlice[0]
		firstWinnerBetAmount := firstWinner.BetAmount()
		pond := int64(0)
		winners := winners[:0]

		for i := 0; i < len(betStatusSlice); i++ {
			eachStatus := betStatusSlice[i]
			if eachStatus.winVal == firstWinner.winVal {
				winners = append(winners, eachStatus.seatID)
			}

			if eachStatus.BetAmount() >= firstWinnerBetAmount {
				eachStatus.betAmount -= firstWinnerBetAmount
				pond += firstWinnerBetAmount
			} else {
				pond += eachStatus.betAmount
				eachStatus.betAmount = 0
			}
		}

		pond /= int64(len(winners))
		for _, wid := range winners {
			seatID2Win[wid] += pond
		}
	}

	return seatID2Win
}

type BetStatus struct {
	seatID    int32
	winVal    uint32
	betAmount int64
}

func NewBetStatus(seatID int32, winVal uint32, betAmount int64) *BetStatus {
	return &BetStatus{
		seatID:    seatID,
		winVal:    winVal,
		betAmount: betAmount,
	}
}

func (b *BetStatus) String() string {
	return fmt.Sprintln("userID", b.seatID, "winVal", b.winVal, "betAmount", b.betAmount)
}

func (b *BetStatus) SeatID() int32 {
	return b.seatID
}

func (b *BetStatus) WinVal() uint32 {
	return b.winVal
}

func (b *BetStatus) BetAmount() int64 {
	return b.betAmount
}

type BetStatusSlice []*BetStatus

func (p BetStatusSlice) Len() int {
	return len(p)
}

func (p BetStatusSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p BetStatusSlice) Less(i, j int) bool {
	if p[i].winVal > p[j].winVal {
		return true
	} else if p[i].winVal < p[j].winVal {
		return false
	}
	return p[i].betAmount < p[j].betAmount
}
