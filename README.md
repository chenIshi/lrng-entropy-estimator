# lrng-entropy-estimator
Verifying the behavior of LRNG entropy estimation.

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