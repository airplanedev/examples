
const express = require('express');
const shortid = require('@team/shared');

const app = express();

app.get('/', (_, res) => {
	res.end(shortid());
});

app.listen(3000, () => {
	console.log('listening on http://localhost:3000');
});