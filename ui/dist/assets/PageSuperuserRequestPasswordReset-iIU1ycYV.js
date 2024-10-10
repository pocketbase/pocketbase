import{S as M,i as T,s as j,F as z,c as L,m as R,t as w,a as y,d as E,b as v,e as m,f as p,g,h as d,j as B,l as N,k as A,n as D,o as k,p as C,q as G,r as F,u as H,v as I,w as h,x as J,y as P,z as S}from"./index-C6KFFkct.js";function K(u){let e,s,n,l,t,r,c,_,i,a,b,f;return l=new G({props:{class:"form-field required",name:"email",$$slots:{default:[Q,({uniqueId:o})=>({5:o}),({uniqueId:o})=>o?32:0]},$$scope:{ctx:u}}}),{c(){e=m("form"),s=m("div"),s.innerHTML='<h4 class="m-b-xs">Forgotten superuser password</h4> <p>Enter the email associated with your account and we’ll send you a recovery link:</p>',n=v(),L(l.$$.fragment),t=v(),r=m("button"),c=m("i"),_=v(),i=m("span"),i.textContent="Send recovery link",p(s,"class","content txt-center m-b-sm"),p(c,"class","ri-mail-send-line"),p(i,"class","txt"),p(r,"type","submit"),p(r,"class","btn btn-lg btn-block"),r.disabled=u[1],F(r,"btn-loading",u[1]),p(e,"class","m-b-base")},m(o,$){g(o,e,$),d(e,s),d(e,n),R(l,e,null),d(e,t),d(e,r),d(r,c),d(r,_),d(r,i),a=!0,b||(f=H(e,"submit",I(u[3])),b=!0)},p(o,$){const q={};$&97&&(q.$$scope={dirty:$,ctx:o}),l.$set(q),(!a||$&2)&&(r.disabled=o[1]),(!a||$&2)&&F(r,"btn-loading",o[1])},i(o){a||(w(l.$$.fragment,o),a=!0)},o(o){y(l.$$.fragment,o),a=!1},d(o){o&&k(e),E(l),b=!1,f()}}}function O(u){let e,s,n,l,t,r,c,_,i;return{c(){e=m("div"),s=m("div"),s.innerHTML='<i class="ri-checkbox-circle-line"></i>',n=v(),l=m("div"),t=m("p"),r=h("Check "),c=m("strong"),_=h(u[0]),i=h(" for the recovery link."),p(s,"class","icon"),p(c,"class","txt-nowrap"),p(l,"class","content"),p(e,"class","alert alert-success")},m(a,b){g(a,e,b),d(e,s),d(e,n),d(e,l),d(l,t),d(t,r),d(t,c),d(c,_),d(t,i)},p(a,b){b&1&&J(_,a[0])},i:P,o:P,d(a){a&&k(e)}}}function Q(u){let e,s,n,l,t,r,c,_;return{c(){e=m("label"),s=h("Email"),l=v(),t=m("input"),p(e,"for",n=u[5]),p(t,"type","email"),p(t,"id",r=u[5]),t.required=!0,t.autofocus=!0},m(i,a){g(i,e,a),d(e,s),g(i,l,a),g(i,t,a),S(t,u[0]),t.focus(),c||(_=H(t,"input",u[4]),c=!0)},p(i,a){a&32&&n!==(n=i[5])&&p(e,"for",n),a&32&&r!==(r=i[5])&&p(t,"id",r),a&1&&t.value!==i[0]&&S(t,i[0])},d(i){i&&(k(e),k(l),k(t)),c=!1,_()}}}function U(u){let e,s,n,l,t,r,c,_;const i=[O,K],a=[];function b(f,o){return f[2]?0:1}return e=b(u),s=a[e]=i[e](u),{c(){s.c(),n=v(),l=m("div"),t=m("a"),t.textContent="Back to login",p(t,"href","/login"),p(t,"class","link-hint"),p(l,"class","content txt-center")},m(f,o){a[e].m(f,o),g(f,n,o),g(f,l,o),d(l,t),r=!0,c||(_=B(N.call(null,t)),c=!0)},p(f,o){let $=e;e=b(f),e===$?a[e].p(f,o):(A(),y(a[$],1,1,()=>{a[$]=null}),D(),s=a[e],s?s.p(f,o):(s=a[e]=i[e](f),s.c()),w(s,1),s.m(n.parentNode,n))},i(f){r||(w(s),r=!0)},o(f){y(s),r=!1},d(f){f&&(k(n),k(l)),a[e].d(f),c=!1,_()}}}function V(u){let e,s;return e=new z({props:{$$slots:{default:[U]},$$scope:{ctx:u}}}),{c(){L(e.$$.fragment)},m(n,l){R(e,n,l),s=!0},p(n,[l]){const t={};l&71&&(t.$$scope={dirty:l,ctx:n}),e.$set(t)},i(n){s||(w(e.$$.fragment,n),s=!0)},o(n){y(e.$$.fragment,n),s=!1},d(n){E(e,n)}}}function W(u,e,s){let n="",l=!1,t=!1;async function r(){if(!l){s(1,l=!0);try{await C.collection("_superusers").requestPasswordReset(n),s(2,t=!0)}catch(_){C.error(_)}s(1,l=!1)}}function c(){n=this.value,s(0,n)}return[n,l,t,r,c]}class Y extends M{constructor(e){super(),T(this,e,W,V,j,{})}}export{Y as default};
