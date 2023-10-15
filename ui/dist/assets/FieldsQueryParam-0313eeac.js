import{S as L,i as S,s as k,N as E,e as s,b as o,w as $,c as F,f as H,g as M,h as e,m as T,y as q,t as N,a as B,o as I,d as J}from"./index-3f9207c2.js";function O(v){let t,i,x,p,g,n,a,h,c,_,r,b,f,y,u,C,m,d;return r=new E({props:{content:`
                    ?fields=*,expand.relField.name
                `}}),{c(){t=s("tr"),i=s("td"),i.textContent="fields",x=o(),p=s("td"),p.innerHTML='<span class="label">String</span>',g=o(),n=s("td"),a=s("p"),h=$(`Comma separated string of the fields to return in the JSON response
            `),c=s("em"),c.textContent="(by default returns all fields)",_=$(`. Ex.:
            `),F(r.$$.fragment),b=o(),f=s("p"),f.innerHTML="<code>*</code> targets all keys from the specific depth level.",y=o(),u=s("p"),u.textContent="In addition, the following field modifiers are also supported:",C=o(),m=s("ul"),m.innerHTML=`<li><code>:excerpt(maxLength, withEllipsis?)</code> <br/>
                Returns a short plain text version of the field string value.
                <br/>
                Ex.:
                <code>?fields=*,description:excerpt(200,true)</code></li>`,H(i,"id","query-page")},m(l,w){M(l,t,w),e(t,i),e(t,x),e(t,p),e(t,g),e(t,n),e(n,a),e(a,h),e(a,c),e(a,_),T(r,a,null),e(n,b),e(n,f),e(n,y),e(n,u),e(n,C),e(n,m),d=!0},p:q,i(l){d||(N(r.$$.fragment,l),d=!0)},o(l){B(r.$$.fragment,l),d=!1},d(l){l&&I(t),J(r)}}}class Q extends L{constructor(t){super(),S(this,t,null,O,k,{})}}export{Q as F};
