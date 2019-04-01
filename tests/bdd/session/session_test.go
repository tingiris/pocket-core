package session_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pokt-network/pocket-core/common"
	"github.com/pokt-network/pocket-core/session"
	"github.com/pokt-network/pocket-core/types"
	"path/filepath"
)

// ************************************************************************************************************
// Milestone: Sessions
//
// Tentative Timeline (4-6 weeks)
//
// Unanswered Questions?
// - Where is the seed data coming from? 99% sure -> *WORLD STATE*
// - How is the state of the session nodes maintained?
// - Is service request validation ddos safe?
// - Forking behavior
// TODO use a diff serialization method for network wide protocol
// TODO deciding on GID format
// ************************************************************************************************************

var _ = Describe("Session", func() {
	
	Describe("Session Creation \\ Computing", func() {
		// Variables representing the seed data
		var devid = []byte(common.SHA256FromString("test"))
		var blockhash = common.SHA256FromString("test")
		var requestedChain = common.Blockchain{Name: "eth", NetID: "1", Version: "1"}
		var requestedChainHash = common.SHA256FromString(fmt.Sprintf("%v", requestedChain))
		absPath, _ := filepath.Abs("../fixtures/xsmallnodepool.json")
		var nodelist = session.FileToNodes(absPath)
		
		Context("Invalid Seed Data", func() {
			
			Context("Parameters are missing or null", func() {
				
				Context("Missing Devid", func() {
					NoDevIDSeed := session.Seed{BlockHash: blockhash, RequestedChain: requestedChainHash, NodeList: nodelist}
					It("should return missing devid error", func() {
						_, err := session.NewSession(NoDevIDSeed)
						Expect(err).To(Equal(session.NoDevID))
					})
				})
			})
			
			Context("Devid is incorrect...", func() {
				
				Context("Devid is incorrect format", func() {
					invalidDevIDSeed := session.Seed{DevID: []byte("invalidtest"), BlockHash: blockhash, RequestedChain: requestedChainHash, NodeList: nodelist}
					It("should return `invalid developer id` error", func() {
						_, err := session.NewSession(invalidDevIDSeed)
						Expect(err).To(Equal(session.InvalidDevIDFormat))
					})
				})
				Context("Devid is not found in world state", func() {
					// TODO pending
				})
			})
			
			Context("Block Hash is incorrect...", func() {
				
				Context("Not a valid block hash format", func() {
					invalidBlockHashFormatSeed := session.Seed{DevID: devid, BlockHash: []byte("invalidtest"), RequestedChain: requestedChainHash, NodeList: nodelist}
					It("should return `invalid block hash` error", func() {
						_, err := session.NewSession(invalidBlockHashFormatSeed)
						Expect(err).To(Equal(session.InvalidBlockHashFormat))
					})
				})
				
				Context("Block hash is expired", func() {
					
					It("should error", func() {
						// TODO pending
					})
				})
			})
			
			Context("Requested Blockchain is invalid...", func() {
				
				Context("No nodes are associated with a blockchain", func() {
					noNodesSeed := session.Seed{DevID: devid, BlockHash: blockhash, RequestedChain: common.SHA256FromString("nosuchchain"), NodeList: nodelist}
					It("should return `invalid blockchain` error", func() {
						_, err := session.NewSession(noNodesSeed)
						Expect(err).To(Equal(session.InsufficientNodes))
					})
				})
			})
		})
		
		Context("Valid Seed Data", func() {
			absPath, _ := filepath.Abs("../fixtures/mediumnodepool.json")
			validSeed := session.NewSeed(devid, absPath, requestedChainHash, blockhash)
			s, err := session.NewSession(validSeed)
			It("should not have returned any error", func() {
				Expect(err).To(BeNil())
			})
			Describe("Generating a valid session", func() {
				It("should generate a session key", func() {
					Expect(s.Key).ToNot(BeNil())
					Expect(len(s.Key)).ToNot(BeZero())
				})
				
				Describe("Node selection", func() {
					
					It("should find the session.NODECOUNT closest nodes to the session key", func() {
						Expect(len(s.Nodes.ValidatorNodes) + len(s.Nodes.ServiceNodes)).To(Equal(session.NODECOUNT))
					})
					
					It("should contain no duplicated nodes", func() {
						check := types.NewSet()
						combined := append(s.Nodes.ServiceNodes, s.Nodes.ValidatorNodes...)
						for _, node := range combined {
							Expect(check.Contains(node.GID)).To(BeFalse())
							check.Add(node.GID)
						}
					})
					
					Describe("Nodes in an evenly distributed fashion", func() {
						
						Context("Small pool of nodes, small number of trials", func() {
							
							It("should result in evenly distributed nodes", func() {
								// TODO heavy compute
							})
						})
						
						Context("Small pool of nodes, large number of trials", func() {
							
							It("should be evenly distributed", func() {
								// TODO heavy compute
							})
						})
						
						Context("Large pool of nodes, small number of trials", func() {
							
							It("should be evenly distributed", func() {
								// TODO heavy compute
							})
						})
						
						Context("Large pool of nodes, large number of trials", func() {
							
							It("should be evenly distributed", func() {
								// TODO heavy compute
							})
						})
					})
				})
				
				Describe("Role assignment", func() {
					
					It("should assign roles to each node", func() {
						Expect(s.Nodes.DelegatedMinter.GID).ToNot(BeEmpty())
						Expect(len(s.Nodes.ValidatorNodes)).To(Equal(session.MAXVALIDATORS))
						Expect(len(s.Nodes.ServiceNodes)).To(Equal(session.MAXSERVICERS))
					})
					
					It("should check the validity of the assigned roles", func() {
						// TODO need blockchain layer
					})
					
					It("should assign roles to nodes proportional to the protocol guidelines", func() {
						Expect(len(s.Nodes.ValidatorNodes) > len(s.Nodes.ServiceNodes))
					})
				})
				
				Describe("Deterministic from the seed data", func() {
					
					Context("2 sessions derived from valid same seed data", func() {
						It("should be = and valid", func() {
							s1, _ := session.NewSession(validSeed)
							s2, _ := session.NewSession(validSeed)
							s3, _ := session.NewSession(validSeed)
							s4, _ := session.NewSession(validSeed)
							Expect(s1).To(Equal(s2))
							Expect(s2).To(Equal(s3))
							Expect(s3).To(Equal(s4))
						})
					})
					
					Context("2 sessions derived from different valid seed data", func() {
						validSeed1 := session.NewSeed(common.SHA256FromString("devid1"), absPath, requestedChainHash, blockhash)
						validSeed2 := session.NewSeed(common.SHA256FromString("devid2"), absPath, requestedChainHash, blockhash)
						It("should be != and valid", func() {
							s1, _ := session.NewSession(validSeed1)
							s2, _ := session.NewSession(validSeed2)
							Expect(s1).ToNot(Equal(s2))
						})
					})
				})
				
				Describe("Expose node info", func() {
					
					Describe("For the developer", func() {
						
						It("should expose the devID", func() {
							Expect(s.DevID).ToNot(BeNil())
							Expect(len(s.DevID)).ToNot(BeZero())
						})
					})
					
					Describe("For the nodes", func() {
						
						It("should expose the nodes host and port", func() {
							for _, node := range s.Nodes.ValidatorNodes {
								Expect(node.IP).ToNot(BeEmpty())
								Expect(node.Port).ToNot(BeEmpty())
							}
							for _, node := range s.Nodes.ServiceNodes {
								Expect(node.IP).ToNot(BeEmpty())
								Expect(node.Port).ToNot(BeEmpty())
							}
						})
						
						It("should expose the unique identifier", func() {
							for _, node := range s.Nodes.ValidatorNodes {
								Expect(node.GID).ToNot(BeEmpty())
							}
							for _, node := range s.Nodes.ServiceNodes {
								Expect(node.GID).ToNot(BeEmpty())
							}
						})
						
						It("should expose the role", func() {
							for _, node := range s.Nodes.ValidatorNodes {
								Expect(node.Role).To(Or(Equal(session.DELEGATEDMINTER), Equal(session.VALIDATE)))
							}
							for _, node := range s.Nodes.ServiceNodes {
								Expect(node.Role).To(Equal(session.SERVICE))
							}
						})
					})
				})
			})
		})
	})
})
