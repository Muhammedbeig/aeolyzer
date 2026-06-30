const fs = require('fs');
const html = fs.readFileSync('C:\\Users\\Muham\\Downloads\\AEOLyzer Agency\\Dashboard.html', 'utf8');

// Find the textarea placeholder or the word 'Ask anything'
const textareaMatch = html.indexOf('Ask anything. Try typing @');
if (textareaMatch !== -1) {
  let divIdx = textareaMatch;
  for (let i = 0; i < 5; i++) {
    divIdx = html.lastIndexOf('<div', divIdx - 1);
  }
  console.log(html.substring(divIdx, textareaMatch + 1500));
} else {
  console.log("Not found");
}
