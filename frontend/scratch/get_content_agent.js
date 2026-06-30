const fs = require('fs');
const html = fs.readFileSync('C:\\Users\\Muham\\Downloads\\AEOLyzer Agency\\Content agent.html', 'utf8');

const askIdx = html.indexOf('Blog Post');
if (askIdx !== -1) {
  let divIdx = askIdx;
  for (let i = 0; i < 5; i++) {
    divIdx = html.lastIndexOf('<div', divIdx - 1);
  }
  
  const snippet = html.substring(divIdx, askIdx + 2000);
  console.log(snippet);
}
