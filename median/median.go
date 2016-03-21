package median

func Median(a [5]int) int {
	if a[0] < a[1] {
		if a[2] < a[3] {
			if a[0] < a[2] {
				if a[1] < a[4] {
					if a[1] < a[2] {
						if a[2] < a[4] {
							return 2
						} else {
							return 4
						}
					} else {
						if a[1] < a[3] {
							return 1
						} else {
							return 3
						}
					}
				} else {
					if a[2] < a[4] {
						if a[3] < a[4] {
							return 3
						} else {
							return 4
						}
					} else {
						if a[1] < a[2] {
							return 1
						} else {
							return 2
						}
					}
				}
			} else {
				if a[3] < a[4] {
					if a[0] < a[3] {
						if a[1] < a[3] {
							return 1
						} else {
							return 3
						}
					} else {
						if a[0] < a[4] {
							return 0
						} else {
							return 4
						}
					}
				} else {
					if a[0] < a[4] {
						if a[1] < a[4] {
							return 1
						} else {
							return 4
						}
					} else {
						if a[0] < a[3] {
							return 0
						} else {
							return 3
						}
					}
				}
			}
		} else {
			if a[0] < a[3] {
				if a[1] < a[4] {
					if a[1] < a[3] {
						if a[3] < a[4] {
							return 3
						} else {
							return 4
						}
					} else {
						if a[1] < a[2] {
							return 1
						} else {
							return 2
						}
					}
				} else {
					if a[3] < a[4] {
						if a[2] < a[4] {
							return 2
						} else {
							return 4
						}
					} else {
						if a[1] < a[3] {
							return 1
						} else {
							return 3
						}
					}
				}
			} else {
				if a[2] < a[4] {
					if a[0] < a[2] {
						if a[1] < a[2] {
							return 1
						} else {
							return 2
						}
					} else {
						if a[0] < a[4] {
							return 0
						} else {
							return 4
						}
					}
				} else {
					if a[0] < a[4] {
						if a[1] < a[4] {
							return 1
						} else {
							return 4
						}
					} else {
						if a[0] < a[2] {
							return 0
						} else {
							return 2
						}
					}
				}
			}
		}
	} else {
		if a[2] < a[3] {
			if a[0] < a[3] {
				if a[2] < a[4] {
					if a[0] < a[4] {
						if a[0] < a[2] {
							return 2
						} else {
							return 0
						}
					} else {
						if a[1] < a[4] {
							return 4
						} else {
							return 1
						}
					}
				} else {
					if a[0] < a[2] {
						if a[0] < a[4] {
							return 4
						} else {
							return 0
						}
					} else {
						if a[1] < a[2] {
							return 2
						} else {
							return 1
						}
					}
				}
			} else {
				if a[1] < a[4] {
					if a[3] < a[4] {
						if a[1] < a[3] {
							return 3
						} else {
							return 1
						}
					} else {
						if a[2] < a[4] {
							return 4
						} else {
							return 2
						}
					}
				} else {
					if a[1] < a[3] {
						if a[1] < a[2] {
							return 2
						} else {
							return 1
						}
					} else {
						if a[3] < a[4] {
							return 4
						} else {
							return 3
						}
					}
				}
			}
		} else {
			if a[0] < a[2] {
				if a[3] < a[4] {
					if a[0] < a[4] {
						if a[0] < a[3] {
							return 3
						} else {
							return 0
						}
					} else {
						if a[1] < a[4] {
							return 4
						} else {
							return 1
						}
					}
				} else {
					if a[0] < a[3] {
						if a[0] < a[4] {
							return 4
						} else {
							return 0
						}
					} else {
						if a[1] < a[3] {
							return 3
						} else {
							return 1
						}
					}
				}
			} else {
				if a[1] < a[4] {
					if a[2] < a[4] {
						if a[1] < a[2] {
							return 2
						} else {
							return 1
						}
					} else {
						if a[3] < a[4] {
							return 4
						} else {
							return 3
						}
					}
				} else {
					if a[1] < a[2] {
						if a[1] < a[3] {
							return 3
						} else {
							return 1
						}
					} else {
						if a[2] < a[4] {
							return 4
						} else {
							return 2
						}
					}
				}
			}
		}
	}
}
