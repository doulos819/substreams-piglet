const express = require('express');
const app = express();
const data = require('./output.json');

app.get('/', (req, res) => {
  res.json(data);
});

const port = 21819;
app.listen(port, () => {
  console.log(`Server running at http://localhost:${port}`);
});

