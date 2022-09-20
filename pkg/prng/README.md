# PRNG

PRNG provides generators to use for testing. 

## SFC32
Code from // SFC32 is from https://simblob.blogspot.com/2022/05/upgrading-prng.html#more .

Usage:

    sfc32 := prng.SFC32(0, 12345, 0, 1)
    for i := 0; i < 10; i++ {
        log.Println(sfc32())
    }

