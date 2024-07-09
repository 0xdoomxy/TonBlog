function toNano(amount) {
    if (!BN.isBN(amount) && !(typeof amount === 'string')) {
        throw new Error('Please pass numbers as strings or BN objects to avoid precision errors.');
    }

    return ethunit.toWei(amount, 'gwei');
}
