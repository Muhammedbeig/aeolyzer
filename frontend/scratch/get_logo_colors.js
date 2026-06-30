const fs = require('fs');
const html = fs.readFileSync('C:\\Users\\Muham\\Downloads\\AEOLyzer Agency\\Dashboard.html', 'utf8');

const aeolyzerMatch = html.indexOf('AEOlyzer');
if (aeolyzerMatch !== -1) {
  let startIdx = html.lastIndexOf('<svg', aeolyzerMatch);
  if (startIdx === -1) startIdx = Math.max(0, aeolyzerMatch - 500);
  console.log(html.substring(startIdx, aeolyzerMatch + 200));
} else {
  console.log("AEOlyzer not found");
}
