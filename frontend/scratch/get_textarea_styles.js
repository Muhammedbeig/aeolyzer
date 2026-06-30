const fs = require('fs');
const html = fs.readFileSync('C:\\Users\\Muham\\Downloads\\AEOLyzer Agency\\Dashboard.html', 'utf8');

const textareaMatch = html.indexOf('Ask anything. Try typing @');
if (textareaMatch !== -1) {
  let divIdx = textareaMatch;
  for (let i = 0; i < 5; i++) {
    divIdx = html.lastIndexOf('<div', divIdx - 1);
  }
  
  const snippet = html.substring(divIdx, textareaMatch + 1500);
  console.log("Contains dark mode classes?", snippet.includes('dark:'));
}
