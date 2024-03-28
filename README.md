# A tiny framework for simple web-sites
I wanted to make a small site for technical blog-posts etc.
I sought up the most contemporary web-stack of 2023-24, and found the following:
[What is the ideal Tech stack to build a website in 2024?](https://dev.to/jakemackie/web-development-in-2024-29d6)

Disappointed to see it was all Vue.js and React - not that there is anything wrong with those approaches, 
I just thought it was a bit much for my tiny site. The comment section however did not disappoint. Full of many D.Y.O. 
suggestions for tiny web-frameworks - especially the 
[Document Markup Library (DML)](https://dml.efpage.de/DML_homepage/index.html) caught my eye.
However I still felt something was missing, so here is a tiny-framework, with custom tags and typescript:

To begin with typescript in a very vanilla form is needed. Either download it from 
[nodejs.org](https://nodejs.org/en/learn/getting-started/how-to-install-nodejs) or as i do on a mac with brew:
```shell
# on MacOs
brew install node
```
Then install typescript globally:
```shell
npm install typescript -g
```
after that the command `tsc <file.ts>` will transpile a single ts-file to plain js, or as I use it `npx tsc`, which 
transpiles from the root of where the `tsconfig.json` is placed.

That is basically it for now.