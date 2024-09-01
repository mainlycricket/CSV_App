"use strict";(self.webpackChunkcsv_app_docs=self.webpackChunkcsv_app_docs||[]).push([[924],{7097:(e,n,t)=>{t.r(n),t.d(n,{assets:()=>d,contentTitle:()=>l,default:()=>h,frontMatter:()=>r,metadata:()=>a,toc:()=>c});var s=t(4848),i=t(8453);const r={sidebar_position:1,sidebar_label:"Getting Started",title:"Getting Started"},l=void 0,a={id:"getting-started",title:"Getting Started",description:"Installing Requirements",source:"@site/docs/getting-started.md",sourceDirName:".",slug:"/getting-started",permalink:"/CSV_App/docs/getting-started",draft:!1,unlisted:!1,editUrl:"https://github.com/mainlycricket/CSV_App_Docs/tree/main/docs/getting-started.md",tags:[],version:"current",sidebarPosition:1,frontMatter:{sidebar_position:1,sidebar_label:"Getting Started",title:"Getting Started"},sidebar:"tutorialSidebar",next:{title:"Schema Generation",permalink:"/CSV_App/docs/db/schema-generation"}},d={},c=[{value:"Installing Requirements",id:"installing-requirements",level:3},{value:"Setting-up the project",id:"setting-up-the-project",level:3},{value:"General Overview",id:"general-overview",level:3}];function o(e){const n={a:"a",code:"code",h3:"h3",li:"li",ol:"ol",p:"p",pre:"pre",ul:"ul",...(0,i.R)(),...e.components};return(0,s.jsxs)(s.Fragment,{children:[(0,s.jsx)(n.h3,{id:"installing-requirements",children:"Installing Requirements"}),"\n",(0,s.jsx)(n.p,{children:"Download and install the following if not already installed"}),"\n",(0,s.jsxs)(n.ul,{children:["\n",(0,s.jsx)(n.li,{children:(0,s.jsx)(n.a,{href:"https://go.dev/dl/",children:"Go"})}),"\n",(0,s.jsx)(n.li,{children:(0,s.jsx)(n.a,{href:"https://www.postgresql.org/download/",children:"PostgreSQL"})}),"\n",(0,s.jsx)(n.li,{children:(0,s.jsx)(n.a,{href:"https://pkg.go.dev/golang.org/x/tools/cmd/goimports",children:"GoImports"})}),"\n"]}),"\n",(0,s.jsx)(n.p,{children:"These commands should execute successfully:"}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-bash",children:"$ go version\n$ psql --version\n$ goimports --help\n"})}),"\n",(0,s.jsx)(n.h3,{id:"setting-up-the-project",children:"Setting-up the project"}),"\n",(0,s.jsxs)(n.ul,{children:["\n",(0,s.jsxs)(n.li,{children:["Clone the ",(0,s.jsx)(n.a,{href:"https://github.com/mainlycricket/CSV_App",children:"project"})]}),"\n",(0,s.jsxs)(n.li,{children:["Delete the existing ",(0,s.jsx)(n.code,{children:"./app"})," directory"]}),"\n",(0,s.jsxs)(n.li,{children:["Remove the files in ",(0,s.jsx)(n.code,{children:"./data"})," directory and copy your CSV files here."]}),"\n"]}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-bash",children:"$ git clone git@github.com:mainlycricket/CSV_App.git\n$ go build .\n$ cd CSV_App\n$ rm -r app\n$ rm data/*\n# copy csv files in ./data\n"})}),"\n",(0,s.jsx)(n.h3,{id:"general-overview",children:"General Overview"}),"\n",(0,s.jsxs)(n.ol,{children:["\n",(0,s.jsxs)(n.li,{children:["Generate schema and examine it i.e. ",(0,s.jsx)(n.code,{children:"data/schema.json"})]}),"\n"]}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-bash",children:"$ ./CSV_App schema\n"})}),"\n",(0,s.jsxs)(n.ol,{start:"2",children:["\n",(0,s.jsx)(n.li,{children:"Generate SQL file and seed the database"}),"\n"]}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-bash",children:'$ ./CSV_App sql\n$ psql -h localhost -U postgres -c \'CREATE DATABASE "DB_Name"\'\n$ psql -h localhost -U postgres -d "DB_Name" -f data/db.sql\n'})}),"\n",(0,s.jsxs)(n.ol,{start:"3",children:["\n",(0,s.jsxs)(n.li,{children:["Examine ",(0,s.jsx)(n.code,{children:"data/appConfig.json"})," and generate app"]}),"\n"]}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-bash",children:"$ ./CSV_App app\n$ cd app && ./setup.sh\n"})}),"\n",(0,s.jsxs)(n.ol,{start:"4",children:["\n",(0,s.jsx)(n.li,{children:"Modify the .env and start the server"}),"\n"]}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-bash",children:"# modify .env, rebuild and start the server\n$ go build . && ./app\n"})}),"\n",(0,s.jsx)(n.p,{children:"Continue reading to understand all the configurations available"})]})}function h(e={}){const{wrapper:n}={...(0,i.R)(),...e.components};return n?(0,s.jsx)(n,{...e,children:(0,s.jsx)(o,{...e})}):o(e)}},8453:(e,n,t)=>{t.d(n,{R:()=>l,x:()=>a});var s=t(6540);const i={},r=s.createContext(i);function l(e){const n=s.useContext(r);return s.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function a(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(i):e.components||i:l(e.components),s.createElement(r.Provider,{value:n},e.children)}}}]);