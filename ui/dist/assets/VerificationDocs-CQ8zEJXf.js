import{S as le,i as ne,s as ie,T as F,e as m,b as T,w as M,f as v,g as b,h as d,x as Y,U as x,V as ye,k as ee,W as Ce,n as te,t as L,a as j,o as h,r as K,u as oe,R as qe,c as G,m as J,d as Z,Q as Ve,X as fe,C as Ae,p as Ie,Y as de}from"./index-DsEcxL-6.js";function ue(s,t,e){const o=s.slice();return o[4]=t[e],o}function me(s,t,e){const o=s.slice();return o[4]=t[e],o}function pe(s,t){let e,o=t[4].code+"",f,c,r,a;function u(){return t[3](t[4])}return{key:s,first:null,c(){e=m("button"),f=M(o),c=T(),v(e,"class","tab-item"),K(e,"active",t[1]===t[4].code),this.first=e},m(g,C){b(g,e,C),d(e,f),d(e,c),r||(a=oe(e,"click",u),r=!0)},p(g,C){t=g,C&4&&o!==(o=t[4].code+"")&&Y(f,o),C&6&&K(e,"active",t[1]===t[4].code)},d(g){g&&h(e),r=!1,a()}}}function _e(s,t){let e,o,f,c;return o=new qe({props:{content:t[4].body}}),{key:s,first:null,c(){e=m("div"),G(o.$$.fragment),f=T(),v(e,"class","tab-item"),K(e,"active",t[1]===t[4].code),this.first=e},m(r,a){b(r,e,a),J(o,e,null),d(e,f),c=!0},p(r,a){t=r;const u={};a&4&&(u.content=t[4].body),o.$set(u),(!c||a&6)&&K(e,"active",t[1]===t[4].code)},i(r){c||(L(o.$$.fragment,r),c=!0)},o(r){j(o.$$.fragment,r),c=!1},d(r){r&&h(e),Z(o)}}}function Pe(s){let t,e,o,f,c,r,a,u=s[0].name+"",g,C,D,P,H,R,B,O,N,q,V,$=[],Q=new Map,U,I,p=[],y=new Map,A,_=F(s[2]);const X=l=>l[4].code;for(let l=0;l<_.length;l+=1){let i=me(s,_,l),n=X(i);Q.set(n,$[l]=pe(n,i))}let E=F(s[2]);const W=l=>l[4].code;for(let l=0;l<E.length;l+=1){let i=ue(s,E,l),n=W(i);y.set(n,p[l]=_e(n,i))}return{c(){t=m("div"),e=m("strong"),e.textContent="POST",o=T(),f=m("div"),c=m("p"),r=M("/api/collections/"),a=m("strong"),g=M(u),C=M("/confirm-verification"),D=T(),P=m("div"),P.textContent="Body Parameters",H=T(),R=m("table"),R.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>token</span></div></td> <td><span class="label">String</span></td> <td>The token from the verification request email.</td></tr></tbody>',B=T(),O=m("div"),O.textContent="Responses",N=T(),q=m("div"),V=m("div");for(let l=0;l<$.length;l+=1)$[l].c();U=T(),I=m("div");for(let l=0;l<p.length;l+=1)p[l].c();v(e,"class","label label-primary"),v(f,"class","content"),v(t,"class","alert alert-success"),v(P,"class","section-title"),v(R,"class","table-compact table-border m-b-base"),v(O,"class","section-title"),v(V,"class","tabs-header compact combined left"),v(I,"class","tabs-content"),v(q,"class","tabs")},m(l,i){b(l,t,i),d(t,e),d(t,o),d(t,f),d(f,c),d(c,r),d(c,a),d(a,g),d(c,C),b(l,D,i),b(l,P,i),b(l,H,i),b(l,R,i),b(l,B,i),b(l,O,i),b(l,N,i),b(l,q,i),d(q,V);for(let n=0;n<$.length;n+=1)$[n]&&$[n].m(V,null);d(q,U),d(q,I);for(let n=0;n<p.length;n+=1)p[n]&&p[n].m(I,null);A=!0},p(l,[i]){(!A||i&1)&&u!==(u=l[0].name+"")&&Y(g,u),i&6&&(_=F(l[2]),$=x($,i,X,1,l,_,Q,V,ye,pe,null,me)),i&6&&(E=F(l[2]),ee(),p=x(p,i,W,1,l,E,y,I,Ce,_e,null,ue),te())},i(l){if(!A){for(let i=0;i<E.length;i+=1)L(p[i]);A=!0}},o(l){for(let i=0;i<p.length;i+=1)j(p[i]);A=!1},d(l){l&&(h(t),h(D),h(P),h(H),h(R),h(B),h(O),h(N),h(q));for(let i=0;i<$.length;i+=1)$[i].d();for(let i=0;i<p.length;i+=1)p[i].d()}}}function Re(s,t,e){let{collection:o}=t,f=204,c=[];const r=a=>e(1,f=a.code);return s.$$set=a=>{"collection"in a&&e(0,o=a.collection)},e(2,c=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[o,f,c,r]}class Be extends le{constructor(t){super(),ne(this,t,Re,Pe,ie,{collection:0})}}function be(s,t,e){const o=s.slice();return o[4]=t[e],o}function he(s,t,e){const o=s.slice();return o[4]=t[e],o}function ve(s,t){let e,o=t[4].code+"",f,c,r,a;function u(){return t[3](t[4])}return{key:s,first:null,c(){e=m("button"),f=M(o),c=T(),v(e,"class","tab-item"),K(e,"active",t[1]===t[4].code),this.first=e},m(g,C){b(g,e,C),d(e,f),d(e,c),r||(a=oe(e,"click",u),r=!0)},p(g,C){t=g,C&4&&o!==(o=t[4].code+"")&&Y(f,o),C&6&&K(e,"active",t[1]===t[4].code)},d(g){g&&h(e),r=!1,a()}}}function ge(s,t){let e,o,f,c;return o=new qe({props:{content:t[4].body}}),{key:s,first:null,c(){e=m("div"),G(o.$$.fragment),f=T(),v(e,"class","tab-item"),K(e,"active",t[1]===t[4].code),this.first=e},m(r,a){b(r,e,a),J(o,e,null),d(e,f),c=!0},p(r,a){t=r;const u={};a&4&&(u.content=t[4].body),o.$set(u),(!c||a&6)&&K(e,"active",t[1]===t[4].code)},i(r){c||(L(o.$$.fragment,r),c=!0)},o(r){j(o.$$.fragment,r),c=!1},d(r){r&&h(e),Z(o)}}}function Oe(s){let t,e,o,f,c,r,a,u=s[0].name+"",g,C,D,P,H,R,B,O,N,q,V,$=[],Q=new Map,U,I,p=[],y=new Map,A,_=F(s[2]);const X=l=>l[4].code;for(let l=0;l<_.length;l+=1){let i=he(s,_,l),n=X(i);Q.set(n,$[l]=ve(n,i))}let E=F(s[2]);const W=l=>l[4].code;for(let l=0;l<E.length;l+=1){let i=be(s,E,l),n=W(i);y.set(n,p[l]=ge(n,i))}return{c(){t=m("div"),e=m("strong"),e.textContent="POST",o=T(),f=m("div"),c=m("p"),r=M("/api/collections/"),a=m("strong"),g=M(u),C=M("/request-verification"),D=T(),P=m("div"),P.textContent="Body Parameters",H=T(),R=m("table"),R.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>email</span></div></td> <td><span class="label">String</span></td> <td>The auth record email address to send the verification request (if exists).</td></tr></tbody>',B=T(),O=m("div"),O.textContent="Responses",N=T(),q=m("div"),V=m("div");for(let l=0;l<$.length;l+=1)$[l].c();U=T(),I=m("div");for(let l=0;l<p.length;l+=1)p[l].c();v(e,"class","label label-primary"),v(f,"class","content"),v(t,"class","alert alert-success"),v(P,"class","section-title"),v(R,"class","table-compact table-border m-b-base"),v(O,"class","section-title"),v(V,"class","tabs-header compact combined left"),v(I,"class","tabs-content"),v(q,"class","tabs")},m(l,i){b(l,t,i),d(t,e),d(t,o),d(t,f),d(f,c),d(c,r),d(c,a),d(a,g),d(c,C),b(l,D,i),b(l,P,i),b(l,H,i),b(l,R,i),b(l,B,i),b(l,O,i),b(l,N,i),b(l,q,i),d(q,V);for(let n=0;n<$.length;n+=1)$[n]&&$[n].m(V,null);d(q,U),d(q,I);for(let n=0;n<p.length;n+=1)p[n]&&p[n].m(I,null);A=!0},p(l,[i]){(!A||i&1)&&u!==(u=l[0].name+"")&&Y(g,u),i&6&&(_=F(l[2]),$=x($,i,X,1,l,_,Q,V,ye,ve,null,he)),i&6&&(E=F(l[2]),ee(),p=x(p,i,W,1,l,E,y,I,Ce,ge,null,be),te())},i(l){if(!A){for(let i=0;i<E.length;i+=1)L(p[i]);A=!0}},o(l){for(let i=0;i<p.length;i+=1)j(p[i]);A=!1},d(l){l&&(h(t),h(D),h(P),h(H),h(R),h(B),h(O),h(N),h(q));for(let i=0;i<$.length;i+=1)$[i].d();for(let i=0;i<p.length;i+=1)p[i].d()}}}function Ee(s,t,e){let{collection:o}=t,f=204,c=[];const r=a=>e(1,f=a.code);return s.$$set=a=>{"collection"in a&&e(0,o=a.collection)},e(2,c=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[o,f,c,r]}class Me extends le{constructor(t){super(),ne(this,t,Ee,Oe,ie,{collection:0})}}function ke(s,t,e){const o=s.slice();return o[5]=t[e],o[7]=e,o}function $e(s,t,e){const o=s.slice();return o[5]=t[e],o[7]=e,o}function we(s){let t,e,o,f,c;function r(){return s[4](s[7])}return{c(){t=m("button"),e=m("div"),e.textContent=`${s[5].title}`,o=T(),v(e,"class","txt"),v(t,"class","tab-item"),K(t,"active",s[1]==s[7])},m(a,u){b(a,t,u),d(t,e),d(t,o),f||(c=oe(t,"click",r),f=!0)},p(a,u){s=a,u&2&&K(t,"active",s[1]==s[7])},d(a){a&&h(t),f=!1,c()}}}function Te(s){let t,e,o,f;var c=s[5].component;function r(a,u){return{props:{collection:a[0]}}}return c&&(e=de(c,r(s))),{c(){t=m("div"),e&&G(e.$$.fragment),o=T(),v(t,"class","tab-item"),K(t,"active",s[1]==s[7])},m(a,u){b(a,t,u),e&&J(e,t,null),d(t,o),f=!0},p(a,u){if(c!==(c=a[5].component)){if(e){ee();const g=e;j(g.$$.fragment,1,0,()=>{Z(g,1)}),te()}c?(e=de(c,r(a)),G(e.$$.fragment),L(e.$$.fragment,1),J(e,t,o)):e=null}else if(c){const g={};u&1&&(g.collection=a[0]),e.$set(g)}(!f||u&2)&&K(t,"active",a[1]==a[7])},i(a){f||(e&&L(e.$$.fragment,a),f=!0)},o(a){e&&j(e.$$.fragment,a),f=!1},d(a){a&&h(t),e&&Z(e)}}}function Ne(s){var E,W,l,i;let t,e,o=s[0].name+"",f,c,r,a,u,g,C,D=s[0].name+"",P,H,R,B,O,N,q,V,$,Q,U,I;B=new Ve({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${s[2]}');

        ...

        await pb.collection('${(E=s[0])==null?void 0:E.name}').requestVerification('test@example.com');

        // ---
        // (optional) in your custom confirmation page:
        // ---

        await pb.collection('${(W=s[0])==null?void 0:W.name}').confirmVerification('VERIFICATION_TOKEN');
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${s[2]}');

        ...

        await pb.collection('${(l=s[0])==null?void 0:l.name}').requestVerification('test@example.com');

        // ---
        // (optional) in your custom confirmation page:
        // ---

        await pb.collection('${(i=s[0])==null?void 0:i.name}').confirmVerification('VERIFICATION_TOKEN');
    `}});let p=F(s[3]),y=[];for(let n=0;n<p.length;n+=1)y[n]=we($e(s,p,n));let A=F(s[3]),_=[];for(let n=0;n<A.length;n+=1)_[n]=Te(ke(s,A,n));const X=n=>j(_[n],1,1,()=>{_[n]=null});return{c(){t=m("h3"),e=M("Account verification ("),f=M(o),c=M(")"),r=T(),a=m("div"),u=m("p"),g=M("Sends "),C=m("strong"),P=M(D),H=M(" account verification request."),R=T(),G(B.$$.fragment),O=T(),N=m("h6"),N.textContent="API details",q=T(),V=m("div"),$=m("div");for(let n=0;n<y.length;n+=1)y[n].c();Q=T(),U=m("div");for(let n=0;n<_.length;n+=1)_[n].c();v(t,"class","m-b-sm"),v(a,"class","content txt-lg m-b-sm"),v(N,"class","m-b-xs"),v($,"class","tabs-header compact"),v(U,"class","tabs-content"),v(V,"class","tabs")},m(n,w){b(n,t,w),d(t,e),d(t,f),d(t,c),b(n,r,w),b(n,a,w),d(a,u),d(u,g),d(u,C),d(C,P),d(u,H),b(n,R,w),J(B,n,w),b(n,O,w),b(n,N,w),b(n,q,w),b(n,V,w),d(V,$);for(let S=0;S<y.length;S+=1)y[S]&&y[S].m($,null);d(V,Q),d(V,U);for(let S=0;S<_.length;S+=1)_[S]&&_[S].m(U,null);I=!0},p(n,[w]){var se,ae,ce,re;(!I||w&1)&&o!==(o=n[0].name+"")&&Y(f,o),(!I||w&1)&&D!==(D=n[0].name+"")&&Y(P,D);const S={};if(w&5&&(S.js=`
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
    `),B.$set(S),w&10){p=F(n[3]);let k;for(k=0;k<p.length;k+=1){const z=$e(n,p,k);y[k]?y[k].p(z,w):(y[k]=we(z),y[k].c(),y[k].m($,null))}for(;k<y.length;k+=1)y[k].d(1);y.length=p.length}if(w&11){A=F(n[3]);let k;for(k=0;k<A.length;k+=1){const z=ke(n,A,k);_[k]?(_[k].p(z,w),L(_[k],1)):(_[k]=Te(z),_[k].c(),L(_[k],1),_[k].m(U,null))}for(ee(),k=A.length;k<_.length;k+=1)X(k);te()}},i(n){if(!I){L(B.$$.fragment,n);for(let w=0;w<A.length;w+=1)L(_[w]);I=!0}},o(n){j(B.$$.fragment,n),_=_.filter(Boolean);for(let w=0;w<_.length;w+=1)j(_[w]);I=!1},d(n){n&&(h(t),h(r),h(a),h(R),h(O),h(N),h(q),h(V)),Z(B,n),fe(y,n),fe(_,n)}}}function Se(s,t,e){let o,{collection:f}=t;const c=[{title:"Request verification",component:Me},{title:"Confirm verification",component:Be}];let r=0;const a=u=>e(1,r=u);return s.$$set=u=>{"collection"in u&&e(0,f=u.collection)},e(2,o=Ae.getApiExampleUrl(Ie.baseURL)),[f,r,o,c,a]}class Fe extends le{constructor(t){super(),ne(this,t,Se,Ne,ie,{collection:0})}}export{Fe as default};
