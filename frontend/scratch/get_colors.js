const fs = require('fs');
const html = fs.readFileSync('C:\\Users\\Muham\\Downloads\\AEOLyzer Agency\\Dashboard.html', 'utf8');

const tailwindConfigMatch = html.match(/tailwind\.config\s*=\s*({[\s\S]*?})/);
if (tailwindConfigMatch) {
  console.log("Found Tailwind Config in HTML:");
  console.log(tailwindConfigMatch[1].substring(0, 1000));
} else {
  console.log("No tailwind config found in HTML.");
}
