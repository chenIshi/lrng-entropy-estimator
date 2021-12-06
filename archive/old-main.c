#include <stdio.h>
#include <stdlib.h>
#include <time.h>

#include "./include/util.h"

#define LIST_SIZE 1000000

#define RAND_UPPER 1000
#define RAND_LOWER 0

#define RANDOM_BLOCK_SIZE 10
#define IDENTICAL_BLOCK_SIZE 1

#define LRNG_DEPTH 3

// Generate a integer list with multiple RANDOM_BLOCK and IDENTICAL_BLOCK 
// RANDOM_BLOCK are composed with randomized int, while IDENTICAL_BLOCK with identical int 
// RANDOM_BLOCK and IDENTICAL_BLOCK are interwined and cascuaded 
// Always start with RANDOM_BLOCK
void block_list_generator(int* buf, int size){
	srand(time(0));
	int identical_num = (RAND_UPPER + RAND_LOWER) / 2;
	int cascuade_block_size = RANDOM_BLOCK_SIZE + IDENTICAL_BLOCK_SIZE;

	for (int i = 0; i < size; i++) {
        if (i % cascuade_block_size < RANDOM_BLOCK_SIZE) {
			buf[i] = (rand() % (RAND_UPPER - RAND_LOWER + 1)) + RAND_LOWER;
		} else {
			buf[i] = identical_num;
		}
	}
}

// Implementation LRNG entropy estimation 
// Notice that according to different LRNG_DEPTH, the first LRNG_DEPTH in entropies 
// will be unavailable due to its innate constriction
void lrng_entropy_estimator(int* input_buf, float* entropies_buf, int size) {
	if (LRNG_DEPTH < 1) {
		fprintf(stderr, "Wrong configuration: LRNG_DEPTH < 1");
		exit(EXIT_FAILURE);
	}
	for (int i = LRNG_DEPTH-1; i<size; i++) {
		// build multiple delta/difference based on the origin input sequences
		for (int j=0; j < LRNG_DEPTH; j++) {
			// first we have to get its relative coefficients
			// i.e. 1 -> 1,1; 2 -> 1,2,1; 3 -> 1,3,3,1; 4 -> 1,4,6,4,1
			// where is similar to finding combination
		}
	}
}

int main () {
	int *block_list;
	block_list = (int *)malloc(sizeof(int)*LIST_SIZE);

	block_list_generator(block_list, LIST_SIZE);
	
	float *entropy_list;
	entropy_list = (float *)malloc(sizeof(float)*LIST_SIZE);

	lrng_entropy_estimator(block_list, entropy_list, LIST_SIZE);
}