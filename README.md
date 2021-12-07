# LRNG Entropy Estimator
Verifying the behavior of LRNG entropy estimation.

In this paper [The Linux Pseudorandom Number Generator Revisited](https://eprint.iacr.org/2012/251.pdf) there is an  abstraction on how LRNG evaluate the "randomness" from its entropy sources. As we can see, LRNG adopts a "differential"-like approach to efficiently come out with an entropy estimation. Such approach is so elegant and simple (compared to other ones with lengthy mathematical induction) that I wonder if it is actually useful. Originally, I try to evaluate its effectiness and furthermore discovering if there is any limitation, i.e. if I can construct non-random number sequences but with a high "entropy" value in LRNG.

However, that's where things get even more complicated. In order to evaluate the experiment result (entropy estimation from LRNG) from the "ground-truth", we have to define "what is our ground-truth?" and "how to evaluate the effectness?".

For example, first how do we construct a "not-that-random" number sequence? There are multiple ways to achieve it: we can simply mix some identical numbers with some random numbers, or we can instead use some not-secure functions to generate 
numbers with different distributions, like *sin()*.

> Such functions are not suitable for RNG since they are easy to predict. But in our case we only need to verify its "entropy", which has only to do with the probabilistic distribution.

Once we choose one (or more) way(s) to generate numbers with different level of randomness, we still have to quantify the randomness for later evaluation. That is, we still have to look for some proper matrices for evaluation.

### Parameters

- **List_SIZE**  
The size of "random" list under evaluation, normally should be set to a large enough value to ease the experiment variance.

- **RAND_UPPER** & **RAND_LOWER**  
The range of availabe random numbers in the evaluation list.

- **RANDOM_BLOCK_SIZE** & **IDENTICAL_BLOCK_SIZE**  
There are several way to construct a "random" list. However, due to our inspection on LRNG algorithm, we assumes that it is heavily rely on the "neighborhood" to determine its corresponding entropy. As a result, we craft a block-like cascuaded structure for sanity check.

> TODO: attach a figure.

- **LRNG_DEPTH**  
By default, the LRNG apply three layer of "difference" onto the original list. We want to further discuss the outcome if we apply more difference layer.


### Reminder
1. `export GOPATH=~/git-repos/lrng-entropy-estimator`