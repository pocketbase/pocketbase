import{S as k,i as v,s as y,F as w,c as x,m as C,t as g,a as $,d as L,G as H,E as M,g as r,o as a,e as u,b as m,f,u as _,y as p}from"./index.994990a4.js";function P(c){let t,s,e,n,i;return{c(){t=u("div"),t.innerHTML=`<div class="icon"><i class="ri-error-warning-line"></i></div> 
            <div class="content txt-bold"><p>Invalid or expired verification token.</p></div>`,s=m(),e=u("button"),e.textContent="Close",f(t,"class","alert alert-danger"),f(e,"type","button"),f(e,"class","btn btn-secondary btn-block")},m(l,o){r(l,t,o),r(l,s,o),r(l,e,o),n||(i=_(e,"click",c[4]),n=!0)},p,d(l){l&&a(t),l&&a(s),l&&a(e),n=!1,i()}}}function S(c){let t,s,e,n,i;return{c(){t=u("div"),t.innerHTML=`<div class="icon"><i class="ri-checkbox-circle-line"></i></div> 
            <div class="content txt-bold"><p>Successfully verified email address.</p></div>`,s=m(),e=u("button"),e.textContent="Close",f(t,"class","alert alert-success"),f(e,"type","button"),f(e,"class","btn btn-secondary btn-block")},m(l,o){r(l,t,o),r(l,s,o),r(l,e,o),n||(i=_(e,"click",c[3]),n=!0)},p,d(l){l&&a(t),l&&a(s),l&&a(e),n=!1,i()}}}function T(c){let t;return{c(){t=u("div"),t.innerHTML='<div class="loader loader-lg"><em>Please wait...</em></div>',f(t,"class","txt-center")},m(s,e){r(s,t,e)},p,d(s){s&&a(t)}}}function F(c){let t;function s(i,l){return i[1]?T:i[0]?S:P}let e=s(c),n=e(c);return{c(){n.c(),t=M()},m(i,l){n.m(i,l),r(i,t,l)},p(i,l){e===(e=s(i))&&n?n.p(i,l):(n.d(1),n=e(i),n&&(n.c(),n.m(t.parentNode,t)))},d(i){n.d(i),i&&a(t)}}}function V(c){let t,s;return t=new w({props:{nobranding:!0,$$slots:{default:[F]},$$scope:{ctx:c}}}),{c(){x(t.$$.fragment)},m(e,n){C(t,e,n),s=!0},p(e,[n]){const i={};n&67&&(i.$$scope={dirty:n,ctx:e}),t.$set(i)},i(e){s||(g(t.$$.fragment,e),s=!0)},o(e){$(t.$$.fragment,e),s=!1},d(e){L(t,e)}}}function q(c,t,s){let{params:e}=t,n=!1,i=!1;l();async function l(){s(1,i=!0);const d=new H("../");try{await d.users.confirmVerification(e==null?void 0:e.token),s(0,n=!0)}catch{s(0,n=!1)}s(1,i=!1)}const o=()=>window.close(),b=()=>window.close();return c.$$set=d=>{"params"in d&&s(2,e=d.params)},[n,i,e,o,b]}class I extends k{constructor(t){super(),v(this,t,q,V,y,{params:2})}}export{I as default};
