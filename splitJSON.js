const fs = require('fs');

function splitJsonAndText(input) {
  let validObjects = [];
  let invalidObjects = [];
  let openBrackets = 0;
  let closeBrackets = 0;
  let i = 0;

  while (i < input.length) {
    if (input[i] === '{') {
      openBrackets++;
    } else if (input[i] === '}') {
      closeBrackets++;
    }

    if (openBrackets === closeBrackets && openBrackets > 0) {
      const jsonStr = input.substring(0, i + 1);

      try {
        const jsonObj = JSON.parse(jsonStr);
        validObjects.push(jsonObj);
      } catch (error) {
        invalidObjects.push(jsonStr);
      }

      input = input.substring(i + 1);
      openBrackets = 0;
      closeBrackets = 0;
      i = 0;
    } else {
      i++;
    }
  }

  return { validObjects, invalidObjects };
}

// Read the file contents
const fileContents = fs.readFileSync('piglet-poly-dai-stream.json', 'utf8');

// Split the input into JSON objects
const { validObjects, invalidObjects } = splitJsonAndText(fileContents);

// Convert the array of valid JSON objects to a JSON string
const validJson = JSON.stringify(validObjects, null, 2);

// Convert the array of invalid JSON objects to a string
const invalidJson = invalidObjects.join('\n');

// Write the output JSON to files
fs.writeFileSync('valid.json', validJson, 'utf8');
fs.writeFileSync('invalid.txt', invalidJson, 'utf8');

console.log('Output JSON files generated successfully.');
