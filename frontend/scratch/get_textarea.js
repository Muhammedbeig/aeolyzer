const fs = require('fs');
const html = fs.readFileSync('C:\\Users\\Muham\\Downloads\\AEOLyzer Agency\\Dashboard.html', 'utf8');

const askIdx = html.indexOf('Ask anything');
if (askIdx !== -1) {
  const svgIdx = html.indexOf('<svg', askIdx);
  if (svgIdx !== -1) {
    // Print a larger chunk to capture the whole bottom toolbar
    console.log(html.substring(svgIdx - 300, svgIdx + 4000));
  }
}
