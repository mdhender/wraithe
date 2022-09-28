# PRNG

PRNG provides generators to use for testing. 

## LCG32
Constants for LCG32 from Open Adventure.

Usage:

    lcg32 := prng.LCG32(0)
    for i := 0; i < 10; i++ {
        log.Println(lcg32())
    }


## SFC32
Code for SFC32 is from https://simblob.blogspot.com/2022/05/upgrading-prng.html#more .

Usage:

    sfc32 := prng.SFC32(0, 12345, 0, 1)
    for i := 0; i < 10; i++ {
        log.Println(sfc32())
    }

