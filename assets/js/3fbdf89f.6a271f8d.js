"use strict";(self.webpackChunkcsv_app_docs=self.webpackChunkcsv_app_docs||[]).push([[817],{2972:(e,n,s)=>{s.r(n),s.d(n,{assets:()=>d,contentTitle:()=>a,default:()=>h,frontMatter:()=>t,metadata:()=>o,toc:()=>c});var r=s(4848),i=s(8453);const t={sidebar_position:3,sidebar_label:"GET Query Options",title:"GET Query Options"},a=void 0,o={id:"app/get-req",title:"GET Query Options",description:"GET all requests can add various query params to customize the response",source:"@site/docs/app/get-req.md",sourceDirName:"app",slug:"/app/get-req",permalink:"/CSV_App/docs/app/get-req",draft:!1,unlisted:!1,editUrl:"https://github.com/mainlycricket/CSV_App_Docs/docs/docs/app/get-req.md",tags:[],version:"current",sidebarPosition:3,frontMatter:{sidebar_position:3,sidebar_label:"GET Query Options",title:"GET Query Options"},sidebar:"tutorialSidebar",previous:{title:"API Routes",permalink:"/CSV_App/docs/app/routes"},next:{title:"Auth",permalink:"/CSV_App/docs/app/auth"}},d={},c=[{value:"Filters",id:"filters",level:3},{value:"Sorting",id:"sorting",level:3},{value:"Pagination",id:"pagination",level:3}];function l(e){const n={a:"a",admonition:"admonition",code:"code",h3:"h3",li:"li",p:"p",pre:"pre",ul:"ul",...(0,i.R)(),...e.components},{Details:s}=n;return s||function(e,n){throw new Error("Expected "+(n?"component":"object")+" `"+e+"` to be defined: you likely forgot to import, pass, or provide it.")}("Details",!0),(0,r.jsxs)(r.Fragment,{children:[(0,r.jsx)(n.p,{children:"GET all requests can add various query params to customize the response"}),"\n",(0,r.jsxs)(s,{children:[(0,r.jsx)("summary",{children:"Sample students GET all response"}),(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-json",children:'{\n  "success": true,\n  "message": "data fetched successfully",\n  "data": {\n    "next": false,\n    "data": [\n      {\n        "Branch_Id": {\n          "Branch_Id": "branch_1",\n          "Branch_Name": "Computer Science"\n        },\n        "Course_Id": {\n          "Course_Id": "course_1",\n          "Course_Name": "B. Tech."\n        },\n        "Student_Father": "Ajay",\n        "Student_Id": 1,\n        "Student_Name": "Tushar",\n        "college_id": {\n          "college_id": "college_1",\n          "college_name": "IIT Delhi"\n        }\n      },\n      {\n        "Branch_Id": {\n          "Branch_Id": "branch_1",\n          "Branch_Name": "Computer Science"\n        },\n        "Course_Id": {\n          "Course_Id": "course_1",\n          "Course_Name": "B. Tech."\n        },\n        "Student_Father": "Nand",\n        "Student_Id": 2,\n        "Student_Name": "Akshay",\n        "college_id": {\n          "college_id": "college_1",\n          "college_name": "IIT Delhi"\n        }\n      }\n    ]\n  }\n}\n'})}),(0,r.jsxs)(n.ul,{children:["\n",(0,r.jsxs)(n.li,{children:["The ",(0,r.jsx)(n.code,{children:"data"})," object contains two major fields: ",(0,r.jsx)(n.code,{children:"next"})," indicates where more data is avaiable - ",(0,r.jsx)(n.code,{children:"data"})," contains the actual data"]}),"\n"]})]}),"\n",(0,r.jsx)(n.h3,{id:"filters",children:"Filters"}),"\n",(0,r.jsxs)(n.ul,{children:["\n",(0,r.jsxs)(n.li,{children:["\n",(0,r.jsx)(n.p,{children:"For data filtering, the keys should have the same name as column names in schema"}),"\n"]}),"\n",(0,r.jsxs)(n.li,{children:["\n",(0,r.jsx)(n.p,{children:"For each column, multiple values can be passed by separating them with a comma"}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{children:"?standard=5,6\n\nThis filters all the rows which have a value of `5` or `6` for the `standard` column\n"})}),"\n"]}),"\n",(0,r.jsxs)(n.li,{children:["\n",(0,r.jsx)(n.p,{children:"Multiple columns are separated by an ampersand (&)"}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{children:"?standard=5,6&date=2024-01-15\n\nThis filters all the rows which have a value of `5` or `6` for the `standard` column and `2024-01-15` value for the `date` column\n"})}),"\n"]}),"\n",(0,r.jsxs)(n.li,{children:["\n",(0,r.jsx)(n.p,{children:"For array fields, if the array contains any of the passed values, the condition is true"}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{children:"?active_on=2024-02-02,2024-06-02\n\nThis filters all the rows in which the `active_on` array contains `2024-02-02` or `2024-06-02` items\n"})}),"\n"]}),"\n",(0,r.jsxs)(n.li,{children:["\n",(0,r.jsx)(n.p,{children:"For string fields including string array fields, values should be passed in separate pairs. This is needed because string fields may contain special symbols"}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{children:"?name=Tushar&name=Akshay\n\nThis filters all the rows which have a value of `Tushar` or `Akshay` for the `name` column\n"})}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{children:"?colors=green&colors=red\n\nThis filters all the rows in the `colors` array contains the `green` or `red` items\n"})}),"\n"]}),"\n"]}),"\n",(0,r.jsx)(n.admonition,{type:"tip",children:(0,r.jsxs)(n.p,{children:[(0,r.jsx)(n.a,{href:"https://www.freecodecamp.org/news/javascript-url-encode-example-how-to-use-encodeuricomponent-and-encodeuri/",children:"Encode the URL"})," after adding the query params"]})}),"\n",(0,r.jsx)(n.h3,{id:"sorting",children:"Sorting"}),"\n",(0,r.jsxs)(n.ul,{children:["\n",(0,r.jsxs)(n.li,{children:["Sorting can be customized using ",(0,r.jsx)(n.code,{children:"__order"})," query param"]}),"\n",(0,r.jsxs)(n.li,{children:["Data is by-default sorted using the primary key field if the ",(0,r.jsx)(n.code,{children:"__order"})," param is not specified"]}),"\n",(0,r.jsx)(n.li,{children:"Data can be sorted based on multiple fields, which should be separated by a comma"}),"\n",(0,r.jsx)(n.li,{children:"Prefix the column name with a hyphen to sort data in the descending order"}),"\n",(0,r.jsx)(n.li,{children:"For foreign key columns, sorting by columns of the foreign table is not supported"}),"\n"]}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{children:"?__order=productName,-price\n\n- This sorts data by `productName` in ascending order and then by `price` in descending order\n"})}),"\n",(0,r.jsx)(n.h3,{id:"pagination",children:"Pagination"}),"\n",(0,r.jsxs)(n.ul,{children:["\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.code,{children:"__page"})," param to specify the page number, by default 1"]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.code,{children:"__limit"})," param to specify the maximum number of records, by default specified in ",(0,r.jsx)(n.a,{href:"./tables-config#pagination",children:"tables config"})]}),"\n"]}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{children:"?__order=studentName&__page=2&__limit=25\n\n- This sorts data by `studentName` in ascending order\n- Skips the first 25 records, and returns the next 25 records (at most)\n"})})]})}function h(e={}){const{wrapper:n}={...(0,i.R)(),...e.components};return n?(0,r.jsx)(n,{...e,children:(0,r.jsx)(l,{...e})}):l(e)}},8453:(e,n,s)=>{s.d(n,{R:()=>a,x:()=>o});var r=s(6540);const i={},t=r.createContext(i);function a(e){const n=r.useContext(t);return r.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function o(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(i):e.components||i:a(e.components),r.createElement(t.Provider,{value:n},e.children)}}}]);