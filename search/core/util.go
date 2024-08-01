package core

// Splits an array of bytes into chunks. Returns the array of chunks, and the
// starting indices of each chunk.
func SplitIntoChunks(data []byte, num_chunks int) ([][]byte, []int) {
	if num_chunks > len(data) { return nil, nil }

	var chunks [][]byte
	var indices []int
	chunk_size := len(data)/num_chunks

	for i := 0; i < len(data); i += chunk_size {
		end := i + chunk_size

		if end > len(data) {
			end = len(data)
		}
		chunks = append(chunks, data[i:end])
		indices = append(indices, i)
	}
	return chunks, indices
}

func FindBytesInArray(source []byte, target []byte) []int {
    var result []int
    for i, _ := range source {
    	// If the target is longer than the remaining amount of bytes, break early.
    	if i+len(target)-1 >= len(source) { 
    		break 
    	}
    	
        for j := 0; j < len(target); j++ {
            if source[i + j] != target[j] { 
            	break 
            }

            if j == len(target)-1 {
            	result = append(result, i)
        	}
        }
    }
    return result
}
