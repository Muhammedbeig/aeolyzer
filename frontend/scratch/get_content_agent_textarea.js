const fs = require('fs');
const html = fs.readFileSync('C:\\Users\\Muham\\Downloads\\AEOLyzer Agency\\Content agent.html', 'utf8');

const placeholderMatch = html.indexOf('Describe what you want to write...');
if (placeholderMatch !== -1) {
  let divIdx = placeholderMatch;
  for (let i = 0; i < 5; i++) {
    divIdx = html.lastIndexOf('<div', divIdx - 1);
  }
  
  console.log(html.substring(divIdx, placeholderMatch + 1500));
}
