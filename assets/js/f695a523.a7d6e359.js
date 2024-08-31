"use strict";(self.webpackChunkcsv_app_docs=self.webpackChunkcsv_app_docs||[]).push([[162],{5915:(e,n,s)=>{s.r(n),s.d(n,{assets:()=>d,contentTitle:()=>l,default:()=>h,frontMatter:()=>o,metadata:()=>r,toc:()=>c});var i=s(4848),t=s(8453);const o={sidebar_position:3,sidebar_label:"Auth",title:"Auth"},l=void 0,r={id:"app/auth",title:"Auth",description:"Enabling Auth",source:"@site/docs/app/auth.md",sourceDirName:"app",slug:"/app/auth",permalink:"/CSV_App/docs/app/auth",draft:!1,unlisted:!1,editUrl:"https://github.com/mainlycricket/CSV_App_Docs/docs/docs/app/auth.md",tags:[],version:"current",sidebarPosition:3,frontMatter:{sidebar_position:3,sidebar_label:"Auth",title:"Auth"},sidebar:"tutorialSidebar",previous:{title:"API Routes",permalink:"/CSV_App/docs/app/routes"},next:{title:"Tables Config",permalink:"/CSV_App/docs/app/tables-config"}},d={},c=[{value:"Enabling Auth",id:"enabling-auth",level:3},{value:"Register Route",id:"register-route",level:3},{value:"Login Route",id:"login-route",level:3},{value:"Logout Route",id:"logout-route",level:3}];function a(e){const n={a:"a",admonition:"admonition",code:"code",em:"em",h3:"h3",li:"li",p:"p",pre:"pre",ul:"ul",...(0,t.R)(),...e.components},{Details:s}=n;return s||function(e,n){throw new Error("Expected "+(n?"component":"object")+" `"+e+"` to be defined: you likely forgot to import, pass, or provide it.")}("Details",!0),(0,i.jsxs)(i.Fragment,{children:[(0,i.jsx)(n.h3,{id:"enabling-auth",children:"Enabling Auth"}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsx)(n.li,{children:"To enable auth in the application, a dedicated CSV file should be used"}),"\n",(0,i.jsx)(n.li,{children:"It is recommended to keep only auth related info in this table"}),"\n",(0,i.jsxs)(n.li,{children:["This table must have:","\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"username"})," column of ",(0,i.jsx)(n.code,{children:"text"})," type as its primary key"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"password"})," column of ",(0,i.jsx)(n.code,{children:"text"})," type with ",(0,i.jsx)(n.code,{children:"Hash"})," enabled"]}),"\n"]}),"\n"]}),"\n",(0,i.jsxs)(n.li,{children:["This table may also have optional fields to enable role/organization based authorization:","\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.code,{children:"role"})," field of ",(0,i.jsx)(n.code,{children:"text"})," type with ",(0,i.jsx)(n.code,{children:"Enums"})," set"]}),"\n",(0,i.jsxs)(n.li,{children:["any number of ",(0,i.jsx)(n.em,{children:"organizational"})," fields - column names can be user-decided but they must be of ",(0,i.jsx)(n.code,{children:"text"})," type"]}),"\n"]}),"\n"]}),"\n",(0,i.jsxs)(n.li,{children:["Checkout the sample ",(0,i.jsx)(n.code,{children:"login.csv"})," ",(0,i.jsx)(n.a,{href:"https://github.com/mainlycricket/CSV_App/blob/main/data/login.csv",children:"here"})]}),"\n"]}),"\n",(0,i.jsx)(n.admonition,{title:"Note",type:"danger",children:(0,i.jsx)(n.p,{children:"The username, password and role must be named as it is"})}),"\n",(0,i.jsx)(n.admonition,{type:"tip",children:(0,i.jsx)(n.p,{children:"The organizational fields can be used to identify organzation, sub-organization, departments, sub-departments and so on"})}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:["The ",(0,i.jsx)(n.code,{children:"authTable"})," field in ",(0,i.jsx)(n.code,{children:"data/appConfig.json"})," should be the name of auth table from ",(0,i.jsx)(n.code,{children:"schema.json"})," which satisfies the above mentioned constraints"]}),"\n",(0,i.jsxs)(n.li,{children:["The ",(0,i.jsx)(n.code,{children:"orgFields"})," array should contain the column names choosen as ",(0,i.jsx)(n.em,{children:"organizational"})," fields in the ",(0,i.jsx)(n.code,{children:"authTable"})]}),"\n"]}),"\n",(0,i.jsx)(n.h3,{id:"register-route",children:"Register Route"}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:["Send POST request to ",(0,i.jsx)(n.code,{children:"/__auth/register"})," with the JSON request body"]}),"\n"]}),"\n",(0,i.jsxs)(s,{children:[(0,i.jsx)("summary",{children:"Sample Register Request Body"}),(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-json",children:'{\n  "username": "john_doe",\n  "password": "secret",\n  "role": "hod",\n  "college_id": "college_1",\n  "course_id": "course_1",\n  "branch_id": "branch_1"\n}\n'})})]}),"\n",(0,i.jsx)(n.h3,{id:"login-route",children:"Login Route"}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:["Send POST request to ",(0,i.jsx)(n.code,{children:"/__auth/login"})," with JSON request body"]}),"\n"]}),"\n",(0,i.jsxs)(s,{children:[(0,i.jsx)("summary",{children:"Sample Login Request Body"}),(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-json",children:'{\n  "username": "john_doe",\n  "password": "secret"\n}\n'})})]}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:["\n",(0,i.jsxs)(n.p,{children:["A JWT cookie named ",(0,i.jsx)(n.code,{children:"access_token"})," with an expiry date of ",(0,i.jsx)(n.code,{children:"30 days"})," is set to the response object, which is automatically sent along with the future requests"]}),"\n"]}),"\n",(0,i.jsxs)(n.li,{children:["\n",(0,i.jsxs)(n.p,{children:["The JWT token contains user info like ",(0,i.jsx)(n.code,{children:"username"}),", ",(0,i.jsx)(n.code,{children:"role"})," etc."]}),"\n"]}),"\n"]}),"\n",(0,i.jsxs)(s,{children:[(0,i.jsx)("summary",{children:"Sample JWT token info"}),(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-json",children:'{\n  "username": "john_doe",\n  "role": "hod",\n  "college_id": "college_1",\n  "course_id": "course_1",\n  "branch_id": "branch_1"\n}\n'})})]}),"\n",(0,i.jsx)(n.h3,{id:"logout-route",children:"Logout Route"}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:["Send GET request to ",(0,i.jsx)(n.code,{children:"/__auth/logout"})]}),"\n"]})]})}function h(e={}){const{wrapper:n}={...(0,t.R)(),...e.components};return n?(0,i.jsx)(n,{...e,children:(0,i.jsx)(a,{...e})}):a(e)}},8453:(e,n,s)=>{s.d(n,{R:()=>l,x:()=>r});var i=s(6540);const t={},o=i.createContext(t);function l(e){const n=i.useContext(o);return i.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function r(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(t):e.components||t:l(e.components),i.createElement(o.Provider,{value:n},e.children)}}}]);