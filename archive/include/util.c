#include <util.h>

int factorial(int upper, int lower) {
    int fact = 1;

    while(upper > lower) {
        fact *= upper;
        upper--;
    }
    return fact;
}

int combination(int n, int r) {
    return factorial(n, r) / factorial(r, 0);
}