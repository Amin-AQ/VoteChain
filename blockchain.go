package main

import (
	"crypto/sha256"
	"fmt"
)

// Vote represents a single vote cast by a voter.
type Vote struct {
	VoterID   int
	Candidate string
}

// Block represents a block in the blockchain, containing multiple votes.
type Block struct {
	PrevHash    string
	CurrentHash string // Add CurrentHash field
	Votes       []Vote
}

// Blockchain is a slice of Block elements.
var Blockchain []Block

// Candidates is a map to store candidate vote counts.
var Candidates map[string]int

// Registered Voter List is a slice of int elements containing registered voter id's
var RegisteredVoterList []int

// RegisterVoter adds a new voter to the system.
func RegisterVoter(voterID int) {
	if voterID <= 0 {
		fmt.Printf("Invalid Voter Id: %d.\n", voterID)
	} else if !isVoterRegistered(voterID) {
		RegisteredVoterList = append(RegisteredVoterList, voterID)
		fmt.Printf("Voter %d registered.\n", voterID)
	} else {
		fmt.Printf("Voter %d already registered.\n", voterID)
	}

}

// CastVote allows a registered voter to cast a vote.
func CastVote(voterID int, candidate string) {
	// Check if the voter is registered
	if !isVoterRegistered(voterID) {
		fmt.Printf("Invalid voter ID %d\n", voterID)
		return
	}

	// Check if the voter has already cast a vote
	if hasVoted(voterID) {
		fmt.Printf("Voter %d has already cast a vote.\n", voterID)
		return
	}

	// Check if the candidate exists
	if !candidateExists(candidate) {
		fmt.Printf("Candidate %s does not exist.\n", candidate)
		return
	}

	// Add the vote to the current block
	vote := Vote{VoterID: voterID, Candidate: candidate}
	lastBlock := Blockchain[len(Blockchain)-1]
	newBlock := Block{
		CurrentHash: calculateHash(lastBlock, vote),
		PrevHash:    lastBlock.CurrentHash,
		Votes:       []Vote{vote},
	}
	Blockchain = append(Blockchain, newBlock)

	// Update candidate vote count
	Candidates[vote.Candidate]++

	fmt.Printf("Vote cast by Voter %d for %s is recorded.\n", voterID, candidate)
}

// calculateHash calculates the hash of a block.
func calculateHash(block Block, vote Vote) string {

	// create a new block including the new vote
	newB := block
	newB.Votes = []Vote{vote}
	newB.PrevHash = block.CurrentHash
	blockData := fmt.Sprintf("%s;%v", newB.PrevHash, newB.Votes[0])

	hash := sha256.New()
	hash.Write([]byte(blockData))

	res := hash.Sum(nil)
	var hashValue string = fmt.Sprintf("%x", res)
	return hashValue
}

// CalculateElectionResults calculates and displays the winner of the election.
func CalculateElectionResults() {
	fmt.Println("\nElection Results:")
	var winner string
	maxVotes := -1

	// Calculating the winner of the election
	for candidateName, voteCount := range Candidates {
		fmt.Printf("%s: %d votes\n", candidateName, voteCount)
		if voteCount > maxVotes {
			winner = candidateName
			maxVotes = voteCount
		} else if voteCount == maxVotes {
			winner = "Tie"
		}
	}

	if winner != "Tie" {
		fmt.Printf("Winner: %s\n", winner)
	} else {
		fmt.Println("Election resulted in a tie.")
	}
}

// function returns true if a candidate with a given name exists else false
func candidateExists(candidateName string) bool {
	_, found := Candidates[candidateName]
	return found
}

// function returns true if voter is registered else false
func isVoterRegistered(id int) bool {
	// checking if voter already registered
	for _, voter := range RegisteredVoterList {
		if voter == id {
			return true
		}
	}
	return false
}

// This function returns true if a voter with given id has already cast a vote, else returns false
func hasVoted(voterId int) bool {
	// Check if the blockchain has any blocks
	if len(Blockchain) == 0 {
		return false
	}

	// searching for voter's vote
	for _, block := range Blockchain {
		// if vote found than return true
		if len(block.Votes) > 0 && block.Votes[0].VoterID == voterId {
			return true
		}
	}
	// vote not found, return false
	return false
}

func main() {
	// Initialize the blockchain with a genesis block.
	genesisBlock := Block{PrevHash: "", CurrentHash: "", Votes: nil}
	Blockchain = append(Blockchain, genesisBlock)

	// Register candidates
	Candidates = make(map[string]int)
	Candidates["Candidate A"] = 0
	Candidates["Candidate B"] = 0

	// Register voters
	for i := 1; i <= 10; i++ {
		RegisterVoter(i)
	}

	// Simulate voting process
	CastVote(1, "Candidate A")
	CastVote(2, "Candidate B")
	CastVote(3, "Candidate A")
	CastVote(3, "Candidate B") // Attempted Double Voting
	CastVote(4, "Candidate B")
	CastVote(5, "Candidate A")
	CastVote(5, "Candidate A")  // Attempted Double Voting
	CastVote(6, "Candidate B")  // Should Print in case of tie
	CastVote(7, "Candidate C")  // Invalid Candidate ID
	CastVote(11, "Candidate B") // Invalid Voter ID

	// Calculate and display election results
	CalculateElectionResults()

	// Display the blockchain
	fmt.Println("\nBlockchain:")
	for i, block := range Blockchain {
		fmt.Printf("Block %d\n", i)
		fmt.Printf("PrevHash: %s\n", block.PrevHash)
		fmt.Printf("CurrentHash: %s\n", block.CurrentHash)
		fmt.Printf("Votes: %v\n\n", block.Votes)
	}
}
