import{S as le,i as ne,s as ie,X as F,h as m,j as y,z as M,k as v,n as b,o as d,H as W,Y as x,Z as Te,E as ee,_ as qe,G as te,t as L,a as U,v as h,l as H,q as oe,W as Ce,c as Y,m as Z,d as Q,V as Ve,$ as fe,J as Ae,p as Ie,a0 as de}from"./index-SKn09NMF.js";function ue(s,t,e){const o=s.slice();return o[4]=t[e],o}function me(s,t,e){const o=s.slice();return o[4]=t[e],o}function pe(s,t){let e,o=t[4].code+"",f,c,r,a;function u(){return t[3](t[4])}return{key:s,first:null,c(){e=m("button"),f=M(o),c=y(),v(e,"class","tab-item"),H(e,"active",t[1]===t[4].code),this.first=e},m(g,q){b(g,e,q),d(e,f),d(e,c),r||(a=oe(e,"click",u),r=!0)},p(g,q){t=g,q&4&&o!==(o=t[4].code+"")&&W(f,o),q&6&&H(e,"active",t[1]===t[4].code)},d(g){g&&h(e),r=!1,a()}}}function _e(s,t){let e,o,f,c;return o=new Ce({props:{content:t[4].body}}),{key:s,first:null,c(){e=m("div"),Y(o.$$.fragment),f=y(),v(e,"class","tab-item"),H(e,"active",t[1]===t[4].code),this.first=e},m(r,a){b(r,e,a),Z(o,e,null),d(e,f),c=!0},p(r,a){t=r;const u={};a&4&&(u.content=t[4].body),o.$set(u),(!c||a&6)&&H(e,"active",t[1]===t[4].code)},i(r){c||(L(o.$$.fragment,r),c=!0)},o(r){U(o.$$.fragment,r),c=!1},d(r){r&&h(e),Q(o)}}}function Pe(s){let t,e,o,f,c,r,a,u=s[0].name+"",g,q,D,P,j,R,B,E,N,C,V,$=[],z=new Map,K,I,p=[],T=new Map,A,_=F(s[2]);const J=l=>l[4].code;for(let l=0;l<_.length;l+=1){let i=me(s,_,l),n=J(i);z.set(n,$[l]=pe(n,i))}let O=F(s[2]);const G=l=>l[4].code;for(let l=0;l<O.length;l+=1){let i=ue(s,O,l),n=G(i);T.set(n,p[l]=_e(n,i))}return{c(){t=m("div"),e=m("strong"),e.textContent="POST",o=y(),f=m("div"),c=m("p"),r=M("/api/collections/"),a=m("strong"),g=M(u),q=M("/confirm-verification"),D=y(),P=m("div"),P.textContent="Body Parameters",j=y(),R=m("table"),R.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>token</span></div></td> <td><span class="label">String</span></td> <td>The token from the verification request email.</td></tr></tbody>',B=y(),E=m("div"),E.textContent="Responses",N=y(),C=m("div"),V=m("div");for(let l=0;l<$.length;l+=1)$[l].c();K=y(),I=m("div");for(let l=0;l<p.length;l+=1)p[l].c();v(e,"class","label label-primary"),v(f,"class","content"),v(t,"class","alert alert-success"),v(P,"class","section-title"),v(R,"class","table-compact table-border m-b-base"),v(E,"class","section-title"),v(V,"class","tabs-header compact combined left"),v(I,"class","tabs-content"),v(C,"class","tabs")},m(l,i){b(l,t,i),d(t,e),d(t,o),d(t,f),d(f,c),d(c,r),d(c,a),d(a,g),d(c,q),b(l,D,i),b(l,P,i),b(l,j,i),b(l,R,i),b(l,B,i),b(l,E,i),b(l,N,i),b(l,C,i),d(C,V);for(let n=0;n<$.length;n+=1)$[n]&&$[n].m(V,null);d(C,K),d(C,I);for(let n=0;n<p.length;n+=1)p[n]&&p[n].m(I,null);A=!0},p(l,[i]){(!A||i&1)&&u!==(u=l[0].name+"")&&W(g,u),i&6&&(_=F(l[2]),$=x($,i,J,1,l,_,z,V,Te,pe,null,me)),i&6&&(O=F(l[2]),ee(),p=x(p,i,G,1,l,O,T,I,qe,_e,null,ue),te())},i(l){if(!A){for(let i=0;i<O.length;i+=1)L(p[i]);A=!0}},o(l){for(let i=0;i<p.length;i+=1)U(p[i]);A=!1},d(l){l&&(h(t),h(D),h(P),h(j),h(R),h(B),h(E),h(N),h(C));for(let i=0;i<$.length;i+=1)$[i].d();for(let i=0;i<p.length;i+=1)p[i].d()}}}function Re(s,t,e){let{collection:o}=t,f=204,c=[];const r=a=>e(1,f=a.code);return s.$$set=a=>{"collection"in a&&e(0,o=a.collection)},e(2,c=[{code:204,body:"null"},{code:400,body:`
                {
                  "code": 400,
                  "message": "An error occurred while validating the submitted data.",
                  "data": {
                    "token": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `}]),[o,f,c,r]}class Be extends le{constructor(t){super(),ne(this,t,Re,Pe,ie,{collection:0})}}function be(s,t,e){const o=s.slice();return o[4]=t[e],o}function he(s,t,e){const o=s.slice();return o[4]=t[e],o}function ve(s,t){let e,o=t[4].code+"",f,c,r,a;function u(){return t[3](t[4])}return{key:s,first:null,c(){e=m("button"),f=M(o),c=y(),v(e,"class","tab-item"),H(e,"active",t[1]===t[4].code),this.first=e},m(g,q){b(g,e,q),d(e,f),d(e,c),r||(a=oe(e,"click",u),r=!0)},p(g,q){t=g,q&4&&o!==(o=t[4].code+"")&&W(f,o),q&6&&H(e,"active",t[1]===t[4].code)},d(g){g&&h(e),r=!1,a()}}}function ge(s,t){let e,o,f,c;return o=new Ce({props:{content:t[4].body}}),{key:s,first:null,c(){e=m("div"),Y(o.$$.fragment),f=y(),v(e,"class","tab-item"),H(e,"active",t[1]===t[4].code),this.first=e},m(r,a){b(r,e,a),Z(o,e,null),d(e,f),c=!0},p(r,a){t=r;const u={};a&4&&(u.content=t[4].body),o.$set(u),(!c||a&6)&&H(e,"active",t[1]===t[4].code)},i(r){c||(L(o.$$.fragment,r),c=!0)},o(r){U(o.$$.fragment,r),c=!1},d(r){r&&h(e),Q(o)}}}function Ee(s){let t,e,o,f,c,r,a,u=s[0].name+"",g,q,D,P,j,R,B,E,N,C,V,$=[],z=new Map,K,I,p=[],T=new Map,A,_=F(s[2]);const J=l=>l[4].code;for(let l=0;l<_.length;l+=1){let i=he(s,_,l),n=J(i);z.set(n,$[l]=ve(n,i))}let O=F(s[2]);const G=l=>l[4].code;for(let l=0;l<O.length;l+=1){let i=be(s,O,l),n=G(i);T.set(n,p[l]=ge(n,i))}return{c(){t=m("div"),e=m("strong"),e.textContent="POST",o=y(),f=m("div"),c=m("p"),r=M("/api/collections/"),a=m("strong"),g=M(u),q=M("/request-verification"),D=y(),P=m("div"),P.textContent="Body Parameters",j=y(),R=m("table"),R.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>email</span></div></td> <td><span class="label">String</span></td> <td>The auth record email address to send the verification request (if exists).</td></tr></tbody>',B=y(),E=m("div"),E.textContent="Responses",N=y(),C=m("div"),V=m("div");for(let l=0;l<$.length;l+=1)$[l].c();K=y(),I=m("div");for(let l=0;l<p.length;l+=1)p[l].c();v(e,"class","label label-primary"),v(f,"class","content"),v(t,"class","alert alert-success"),v(P,"class","section-title"),v(R,"class","table-compact table-border m-b-base"),v(E,"class","section-title"),v(V,"class","tabs-header compact combined left"),v(I,"class","tabs-content"),v(C,"class","tabs")},m(l,i){b(l,t,i),d(t,e),d(t,o),d(t,f),d(f,c),d(c,r),d(c,a),d(a,g),d(c,q),b(l,D,i),b(l,P,i),b(l,j,i),b(l,R,i),b(l,B,i),b(l,E,i),b(l,N,i),b(l,C,i),d(C,V);for(let n=0;n<$.length;n+=1)$[n]&&$[n].m(V,null);d(C,K),d(C,I);for(let n=0;n<p.length;n+=1)p[n]&&p[n].m(I,null);A=!0},p(l,[i]){(!A||i&1)&&u!==(u=l[0].name+"")&&W(g,u),i&6&&(_=F(l[2]),$=x($,i,J,1,l,_,z,V,Te,ve,null,he)),i&6&&(O=F(l[2]),ee(),p=x(p,i,G,1,l,O,T,I,qe,ge,null,be),te())},i(l){if(!A){for(let i=0;i<O.length;i+=1)L(p[i]);A=!0}},o(l){for(let i=0;i<p.length;i+=1)U(p[i]);A=!1},d(l){l&&(h(t),h(D),h(P),h(j),h(R),h(B),h(E),h(N),h(C));for(let i=0;i<$.length;i+=1)$[i].d();for(let i=0;i<p.length;i+=1)p[i].d()}}}function Oe(s,t,e){let{collection:o}=t,f=204,c=[];const r=a=>e(1,f=a.code);return s.$$set=a=>{"collection"in a&&e(0,o=a.collection)},e(2,c=[{code:204,body:"null"},{code:400,body:`
                {
                  "code": 400,
                  "message": "An error occurred while validating the submitted data.",
                  "data": {
                    "email": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `}]),[o,f,c,r]}class Me extends le{constructor(t){super(),ne(this,t,Oe,Ee,ie,{collection:0})}}function ke(s,t,e){const o=s.slice();return o[5]=t[e],o[7]=e,o}function $e(s,t,e){const o=s.slice();return o[5]=t[e],o[7]=e,o}function we(s){let t,e,o,f,c;function r(){return s[4](s[7])}return{c(){t=m("button"),e=m("div"),e.textContent=`${s[5].title}`,o=y(),v(e,"class","txt"),v(t,"class","tab-item"),H(t,"active",s[1]==s[7])},m(a,u){b(a,t,u),d(t,e),d(t,o),f||(c=oe(t,"click",r),f=!0)},p(a,u){s=a,u&2&&H(t,"active",s[1]==s[7])},d(a){a&&h(t),f=!1,c()}}}function ye(s){let t,e,o,f;var c=s[5].component;function r(a,u){return{props:{collection:a[0]}}}return c&&(e=de(c,r(s))),{c(){t=m("div"),e&&Y(e.$$.fragment),o=y(),v(t,"class","tab-item"),H(t,"active",s[1]==s[7])},m(a,u){b(a,t,u),e&&Z(e,t,null),d(t,o),f=!0},p(a,u){if(c!==(c=a[5].component)){if(e){ee();const g=e;U(g.$$.fragment,1,0,()=>{Q(g,1)}),te()}c?(e=de(c,r(a)),Y(e.$$.fragment),L(e.$$.fragment,1),Z(e,t,o)):e=null}else if(c){const g={};u&1&&(g.collection=a[0]),e.$set(g)}(!f||u&2)&&H(t,"active",a[1]==a[7])},i(a){f||(e&&L(e.$$.fragment,a),f=!0)},o(a){e&&U(e.$$.fragment,a),f=!1},d(a){a&&h(t),e&&Q(e)}}}function Ne(s){var O,G,l,i;let t,e,o=s[0].name+"",f,c,r,a,u,g,q,D=s[0].name+"",P,j,R,B,E,N,C,V,$,z,K,I;B=new Ve({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${s[2]}');

        ...

        await pb.collection('${(O=s[0])==null?void 0:O.name}').requestVerification('test@example.com');

        // ---
        // (optional) in your custom confirmation page:
        // ---

        await pb.collection('${(G=s[0])==null?void 0:G.name}').confirmVerification('VERIFICATION_TOKEN');
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${s[2]}');

        ...

        await pb.collection('${(l=s[0])==null?void 0:l.name}').requestVerification('test@example.com');

        // ---
        // (optional) in your custom confirmation page:
        // ---

        await pb.collection('${(i=s[0])==null?void 0:i.name}').confirmVerification('VERIFICATION_TOKEN');
    `}});let p=F(s[3]),T=[];for(let n=0;n<p.length;n+=1)T[n]=we($e(s,p,n));let A=F(s[3]),_=[];for(let n=0;n<A.length;n+=1)_[n]=ye(ke(s,A,n));const J=n=>U(_[n],1,1,()=>{_[n]=null});return{c(){t=m("h3"),e=M("Account verification ("),f=M(o),c=M(")"),r=y(),a=m("div"),u=m("p"),g=M("Sends "),q=m("strong"),P=M(D),j=M(" account verification request."),R=y(),Y(B.$$.fragment),E=y(),N=m("h6"),N.textContent="API details",C=y(),V=m("div"),$=m("div");for(let n=0;n<T.length;n+=1)T[n].c();z=y(),K=m("div");for(let n=0;n<_.length;n+=1)_[n].c();v(t,"class","m-b-sm"),v(a,"class","content txt-lg m-b-sm"),v(N,"class","m-b-xs"),v($,"class","tabs-header compact"),v(K,"class","tabs-content"),v(V,"class","tabs")},m(n,w){b(n,t,w),d(t,e),d(t,f),d(t,c),b(n,r,w),b(n,a,w),d(a,u),d(u,g),d(u,q),d(q,P),d(u,j),b(n,R,w),Z(B,n,w),b(n,E,w),b(n,N,w),b(n,C,w),b(n,V,w),d(V,$);for(let S=0;S<T.length;S+=1)T[S]&&T[S].m($,null);d(V,z),d(V,K);for(let S=0;S<_.length;S+=1)_[S]&&_[S].m(K,null);I=!0},p(n,[w]){var se,ae,ce,re;(!I||w&1)&&o!==(o=n[0].name+"")&&W(f,o),(!I||w&1)&&D!==(D=n[0].name+"")&&W(P,D);const S={};if(w&5&&(S.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${n[2]}');

        ...

        await pb.collection('${(se=n[0])==null?void 0:se.name}').requestVerification('test@example.com');

        // ---
        // (optional) in your custom confirmation page:
        // ---

        await pb.collection('${(ae=n[0])==null?void 0:ae.name}').confirmVerification('VERIFICATION_TOKEN');
    `),w&5&&(S.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${n[2]}');

        ...

        await pb.collection('${(ce=n[0])==null?void 0:ce.name}').requestVerification('test@example.com');

        // ---
        // (optional) in your custom confirmation page:
        // ---

        await pb.collection('${(re=n[0])==null?void 0:re.name}').confirmVerification('VERIFICATION_TOKEN');
    `),B.$set(S),w&10){p=F(n[3]);let k;for(k=0;k<p.length;k+=1){const X=$e(n,p,k);T[k]?T[k].p(X,w):(T[k]=we(X),T[k].c(),T[k].m($,null))}for(;k<T.length;k+=1)T[k].d(1);T.length=p.length}if(w&11){A=F(n[3]);let k;for(k=0;k<A.length;k+=1){const X=ke(n,A,k);_[k]?(_[k].p(X,w),L(_[k],1)):(_[k]=ye(X),_[k].c(),L(_[k],1),_[k].m(K,null))}for(ee(),k=A.length;k<_.length;k+=1)J(k);te()}},i(n){if(!I){L(B.$$.fragment,n);for(let w=0;w<A.length;w+=1)L(_[w]);I=!0}},o(n){U(B.$$.fragment,n),_=_.filter(Boolean);for(let w=0;w<_.length;w+=1)U(_[w]);I=!1},d(n){n&&(h(t),h(r),h(a),h(R),h(E),h(N),h(C),h(V)),Q(B,n),fe(T,n),fe(_,n)}}}function Se(s,t,e){let o,{collection:f}=t;const c=[{title:"Request verification",component:Me},{title:"Confirm verification",component:Be}];let r=0;const a=u=>e(1,r=u);return s.$$set=u=>{"collection"in u&&e(0,f=u.collection)},e(2,o=Ae.getApiExampleUrl(Ie.baseURL)),[f,r,o,c,a]}class Fe extends le{constructor(t){super(),ne(this,t,Se,Ne,ie,{collection:0})}}export{Fe as default};
