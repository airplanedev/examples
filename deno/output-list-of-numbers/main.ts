[...Array(5)].forEach((v, i) => {
    setTimeout(() => {
        console.log(`airplane_output: ${i}`);
    }, i * 500);
});