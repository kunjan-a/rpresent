package main

// Imported from https://code.google.com/p/go/source/browse/present/static/styles.css?repo=talks
const stylesCSS = `
/* Framework */

html {
  height: 100%;
}

body {
  margin: 0;
  padding: 0;

  display: block !important;

  height: 100%;
  min-height: 740px;

  overflow-x: hidden;
  overflow-y: auto;

  background: rgb(215, 215, 215);
  background: -o-radial-gradient(rgb(240, 240, 240), rgb(190, 190, 190));
  background: -moz-radial-gradient(rgb(240, 240, 240), rgb(190, 190, 190));
  background: -webkit-radial-gradient(rgb(240, 240, 240), rgb(190, 190, 190));
  background: -webkit-gradient(radial, 50% 50%, 0, 50% 50%, 500, from(rgb(240, 240, 240)), to(rgb(190, 190, 190)));

  -webkit-font-smoothing: antialiased;
}

.helpLink {
  display: block;
  position: absolute;
  right: 50px;
  font-size: 8pt;
}

.slides {
  width: 100%;
  height: 100%;
  left: 0;
  top: 0;

  position: absolute;

  -webkit-transform: translate3d(0, 0, 0);
}

.slides > article {
  display: block;

  position: absolute;
  overflow: hidden;

  width: 900px;
  height: 700px;

  left: 50%;
  top: 50%;

  margin-left: -450px;
  margin-top: -350px;

  padding: 40px 60px;

  box-sizing: border-box;
  -o-box-sizing: border-box;
  -moz-box-sizing: border-box;
  -webkit-box-sizing: border-box;

  border-radius: 10px;
  -o-border-radius: 10px;
  -moz-border-radius: 10px;
  -webkit-border-radius: 10px;

  background-color: white;

  border: 1px solid rgba(0, 0, 0, .3);

  transition: transform .3s ease-out;
  -o-transition: -o-transform .3s ease-out;
  -moz-transition: -moz-transform .3s ease-out;
  -webkit-transition: -webkit-transform .3s ease-out;
}
.slides.layout-widescreen > article {
  margin-left: -550px;
  width: 1100px;
}
.slides.layout-faux-widescreen > article {
  margin-left: -550px;
  width: 1100px;

  padding: 40px 160px;
}

.slides.layout-widescreen > article:not(.nobackground):not(.biglogo),
.slides.layout-faux-widescreen > article:not(.nobackground):not(.biglogo) {
  background-position-x: 0, 840px;
}

/* Clickable/tappable areas */

.slide-area {
  z-index: 1000;

  position: absolute;
  left: 0;
  top: 0;
  width: 150px;
  height: 700px;

  left: 50%;
  top: 50%;

  cursor: pointer;
  margin-top: -350px;

  tap-highlight-color: transparent;
  -o-tap-highlight-color: transparent;
  -moz-tap-highlight-color: transparent;
  -webkit-tap-highlight-color: transparent;
}
#prev-slide-area {
  margin-left: -550px;
}
#next-slide-area {
  margin-left: 400px;
}
.slides.layout-widescreen #prev-slide-area,
.slides.layout-faux-widescreen #prev-slide-area {
  margin-left: -650px;
}
.slides.layout-widescreen #next-slide-area,
.slides.layout-faux-widescreen #next-slide-area {
  margin-left: 500px;
}

/* Slides */

.slides > article {
  display: none;
}
.slides > article.far-past {
  display: block;
  transform: translate(-2040px);
  -o-transform: translate(-2040px);
  -moz-transform: translate(-2040px);
  -webkit-transform: translate3d(-2040px, 0, 0);
}
.slides > article.past {
  display: block;
  transform: translate(-1020px);
  -o-transform: translate(-1020px);
  -moz-transform: translate(-1020px);
  -webkit-transform: translate3d(-1020px, 0, 0);
}
.slides > article.current {
  display: block;
  transform: translate(0);
  -o-transform: translate(0);
  -moz-transform: translate(0);
  -webkit-transform: translate3d(0, 0, 0);
}
.slides > article.next {
  display: block;
  transform: translate(1020px);
  -o-transform: translate(1020px);
  -moz-transform: translate(1020px);
  -webkit-transform: translate3d(1020px, 0, 0);
}
.slides > article.far-next {
  display: block;
  transform: translate(2040px);
  -o-transform: translate(2040px);
  -moz-transform: translate(2040px);
  -webkit-transform: translate3d(2040px, 0, 0);
}

.slides.layout-widescreen > article.far-past,
.slides.layout-faux-widescreen > article.far-past {
  display: block;
  transform: translate(-2260px);
  -o-transform: translate(-2260px);
  -moz-transform: translate(-2260px);
  -webkit-transform: translate3d(-2260px, 0, 0);
}
.slides.layout-widescreen > article.past,
.slides.layout-faux-widescreen > article.past {
  display: block;
  transform: translate(-1130px);
  -o-transform: translate(-1130px);
  -moz-transform: translate(-1130px);
  -webkit-transform: translate3d(-1130px, 0, 0);
}
.slides.layout-widescreen > article.current,
.slides.layout-faux-widescreen > article.current {
  display: block;
  transform: translate(0);
  -o-transform: translate(0);
  -moz-transform: translate(0);
  -webkit-transform: translate3d(0, 0, 0);
}
.slides.layout-widescreen > article.next,
.slides.layout-faux-widescreen > article.next {
  display: block;
  transform: translate(1130px);
  -o-transform: translate(1130px);
  -moz-transform: translate(1130px);
  -webkit-transform: translate3d(1130px, 0, 0);
}
.slides.layout-widescreen > article.far-next,
.slides.layout-faux-widescreen > article.far-next {
  display: block;
  transform: translate(2260px);
  -o-transform: translate(2260px);
  -moz-transform: translate(2260px);
  -webkit-transform: translate3d(2260px, 0, 0);
}

/* Styles for slides */

.slides > article {
  font-family: 'Open Sans', Arial, sans-serif;

  color: black;
  text-shadow: 0 1px 1px rgba(0, 0, 0, .1);

  font-size: 26px;
  line-height: 36px;

  letter-spacing: -1px;
}

b {
  font-weight: 600;
}

a {
  color: rgb(0, 102, 204);
  text-decoration: none;
}
a:visited {
  color: rgba(0, 102, 204, .75);
}
a:hover {
  color: black;
}

p {
  margin: 0;
  padding: 0;

  margin-top: 20px;
}
p:first-child {
  margin-top: 0;
}

h1 {
  font-size: 60px;
  line-height: 60px;

  padding: 0;
  margin: 0;
  margin-top: 200px;
  margin-bottom: 5px;
  padding-right: 40px;

  font-weight: 600;

  letter-spacing: -3px;

  color: rgb(51, 51, 51);
}

h2 {
  font-size: 45px;
  line-height: 45px;

  position: absolute;
  bottom: 150px;

  padding: 0;
  margin: 0;
  padding-right: 40px;

  font-weight: 600;

  letter-spacing: -2px;

  color: rgb(51, 51, 51);
}

h3 {
  font-size: 30px;
  line-height: 36px;

  padding: 0;
  margin: 0;
  padding-right: 40px;

  font-weight: 600;

  letter-spacing: -1px;

  color: rgb(51, 51, 51);
}

ul {
  margin: 0;
  padding: 0;
  margin-top: 20px;
  margin-left: 1.5em;
}
li {
  padding: 0;
  margin: 0 0 .5em 0;
}

div.code {
  padding: 5px 10px;
  margin-top: 20px;
  margin-bottom: 20px;
  overflow: hidden;

  background: rgb(240, 240, 240);
  border: 1px solid rgb(224, 224, 224);
}
pre {
  margin: 0;
  padding: 0;

  font-family: 'Droid Sans Mono', 'Courier New', monospace;
  font-size: 18px;
  line-height: 24px;
  letter-spacing: -1px;

  color: black;
}

code {
  font-size: 95%;
  font-family: 'Droid Sans Mono', 'Courier New', monospace;

  color: black;
}

article > .image {
  	text-align: center;
    margin-top: 40px;
}

table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 40px;
}
th {
  font-weight: 600;
  text-align: left;
}
td,
th {
  border: 1px solid rgb(224, 224, 224);
  padding: 5px 10px;
  vertical-align: top;
}

p.link {
  margin-left: 20px;
}

/* Code */
div.code {
  outline: 0px solid transparent;
}
div.playground {
  position: relative;
}
div.output {
  position: absolute;
  left: 50%;
  top: 50%;
  right: 40px;
  bottom: 40px;
  background: #202020;
  padding: 5px 10px;
  z-index: 2;

  border-radius: 10px;
  -o-border-radius: 10px;
  -moz-border-radius: 10px;
  -webkit-border-radius: 10px;

}
div.output pre {
  margin: 0;
  padding: 0;
  background: none;
  border: none;
  width: 100%;
  height: 100%;
  overflow: auto;
}
div.output .stdout, div.output pre {
  color: #e6e6e6;
}
div.output .stderr, div.output .error {
  color: rgb(244, 74, 63);
}
div.output .system, div.output .exit {
  color: rgb(255, 209, 77)
}
.buttons {
  position: relative;
  float: right;
  top: -60px;
  right: 10px;
}
div.output .buttons {
  position: absolute;
  float: none;
  top: auto;
  right: 5px;
  bottom: 5px;
}

/* Presenter details */
.presenter {
	margin-top: 20px;
}
.presenter p,
.presenter .link {
	margin: 0;
	font-size: 28px;
	line-height: 1.2em;
}

/* Output resize details */
.ui-resizable-handle {
  position: absolute;
}
.ui-resizable-n {
  cursor: n-resize;
  height: 7px;
  width: 100%;
  top: -5px;
  left: 0;
}
.ui-resizable-w {
  cursor: w-resize;
  width: 7px;
  left: -5px;
  top: 0;
  height: 100%;
}
.ui-resizable-nw {
  cursor: nw-resize;
  width: 9px;
  height: 9px;
  left: -5px;
  top: -5px;
}`
