package state

import (
	"fmt"
	)

	// Node represents any MPT node type (Branch, Extension, Leaf)
	type Node interface{}

	// LeafNode stores [encoded_path, value]
	type LeafNode struct {
		Path  []byte
			Value []byte
			}

			// ExtensionNode stores [encoded_path, child_node]
			type ExtensionNode struct {
				Path  []byte
					Child Node
					}

					// BranchNode stores 16 children (0-f) + 1 optional value slot
					type BranchNode struct {
						ChildrenNode
							Value    []byte
							}

							// BytesToNibbles converts a byte slice into a slice of 4-bit nibbles.
							func BytesToNibbles(b []byte) []byte {
								nibbles := make([]byte, len(b)*2)
									for i, byteVal := range b {
											nibbles[i*2] = byteVal >> 4
													nibbles[i*2+1] = byteVal & 0x0f
														}
															return nibbles
															}

															// HexPrefixEncode encodes a slice of nibbles with a flag prefix (Yellow Paper Appendix C).
															// isLeaf indicates if the path terminates a key (Leaf vs Extension node).
															func HexPrefixEncode(nibbles []byte, isLeaf bool) []byte {
																var flag byte
																	if isLeaf {
																			flag = 2
																				} else {
																						flag = 0
																							}

																								oddLen := len(nibbles)%2 != 0
																									if oddLen {
																											flag |= 1
																												}

																													var result []byte
																														if oddLen {
																																// Odd length: High nibble is flag, Low nibble is first key nibble
																																		firstByte := (flag << 4) | (nibbles & 0x0f)
																																				result = append(result, firstByte)
																																						nibbles = nibbles[1:]
																																							} else {
																																									// Even length: High nibble is flag, Low nibble is 0x0 padding
																																											firstByte := flag << 4
																																													result = append(result, firstByte)
																																														}

																																															// Pack remaining nibbles pairwise into bytes
																																																for i := 0; i < len(nibbles); i += 2 {
																																																		b := (nibbles[i] << 4) | (nibbles[i+1] & 0x0f)
																																																				result = append(result, b)
																																																					}

																																																						return result
																																																						}

																																																						// HexPrefixDecode decodes Hex-Prefix encoded bytes back to nibbles and returns isLeaf.
																																																						func HexPrefixDecode(encoded []byte) ([]byte, bool, error) {
																																																							if len(encoded) == 0 {
																																																									return nil, false, fmt.Errorf("empty HP encoded bytes")
																																																										}

																																																											firstByte := encoded
																																																												flag := firstByte >> 4
																																																													isLeaf := (flag & 2) != 0
																																																														isOdd := (flag & 1) != 0

																																																															var nibbles []byte

																																																																if isOdd {
																																																																		nibbles = append(nibbles, firstByte&0x0f)
																																																																			}

																																																																				for _, b := range encoded[1:] {
																																																																						nibbles = append(nibbles, b>>4, b&0x0f)
																																																																							}

																																																																								return nibbles, isLeaf, nil
																																																																								}
																																																																								