import{S as le,i as ne,s as ie,X as F,h as b,t as j,a as U,I as Y,Z as x,_ as Te,C as ee,$ as Ce,D as te,l as h,n as u,u as m,v as y,A as M,w as v,k as K,o as oe,W as qe,d as z,m as G,c as Q,V as Ve,Y as fe,J as Ae,p as Ie,a0 as ue}from"./index-CkK5VYgS.js";function de(s,t,e){const o=s.slice();return o[4]=t[e],o}function me(s,t,e){const o=s.slice();return o[4]=t[e],o}function pe(s,t){let e,o=t[4].code+"",f,c,r,a;function d(){return t[3](t[4])}return{key:s,first:null,c(){e=m("button"),f=M(o),c=y(),v(e,"class","tab-item"),K(e,"active",t[1]===t[4].code),this.first=e},m(g,C){h(g,e,C),u(e,f),u(e,c),r||(a=oe(e,"click",d),r=!0)},p(g,C){t=g,C&4&&o!==(o=t[4].code+"")&&Y(f,o),C&6&&K(e,"active",t[1]===t[4].code)},d(g){g&&b(e),r=!1,a()}}}function _e(s,t){let e,o,f,c;return o=new qe({props:{content:t[4].body}}),{key:s,first:null,c(){e=m("div"),Q(o.$$.fragment),f=y(),v(e,"class","tab-item"),K(e,"active",t[1]===t[4].code),this.first=e},m(r,a){h(r,e,a),G(o,e,null),u(e,f),c=!0},p(r,a){t=r;const d={};a&4&&(d.content=t[4].body),o.$set(d),(!c||a&6)&&K(e,"active",t[1]===t[4].code)},i(r){c||(U(o.$$.fragment,r),c=!0)},o(r){j(o.$$.fragment,r),c=!1},d(r){r&&b(e),z(o)}}}function Pe(s){let t,e,o,f,c,r,a,d=s[0].name+"",g,C,D,P,L,R,B,O,N,q,V,$=[],J=new Map,H,I,p=[],T=new Map,A,_=F(s[2]);const X=l=>l[4].code;for(let l=0;l<_.length;l+=1){let i=me(s,_,l),n=X(i);J.set(n,$[l]=pe(n,i))}let E=F(s[2]);const W=l=>l[4].code;for(let l=0;l<E.length;l+=1){let i=de(s,E,l),n=W(i);T.set(n,p[l]=_e(n,i))}return{c(){t=m("div"),e=m("strong"),e.textContent="POST",o=y(),f=m("div"),c=m("p"),r=M("/api/collections/"),a=m("strong"),g=M(d),C=M("/confirm-verification"),D=y(),P=m("div"),P.textContent="Body Parameters",L=y(),R=m("table"),R.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>token</span></div></td> <td><span class="label">String</span></td> <td>The token from the verification request email.</td></tr></tbody>',B=y(),O=m("div"),O.textContent="Responses",N=y(),q=m("div"),V=m("div");for(let l=0;l<$.length;l+=1)$[l].c();H=y(),I=m("div");for(let l=0;l<p.length;l+=1)p[l].c();v(e,"class","label label-primary"),v(f,"class","content"),v(t,"class","alert alert-success"),v(P,"class","section-title"),v(R,"class","table-compact table-border m-b-base"),v(O,"class","section-title"),v(V,"class","tabs-header compact combined left"),v(I,"class","tabs-content"),v(q,"class","tabs")},m(l,i){h(l,t,i),u(t,e),u(t,o),u(t,f),u(f,c),u(c,r),u(c,a),u(a,g),u(c,C),h(l,D,i),h(l,P,i),h(l,L,i),h(l,R,i),h(l,B,i),h(l,O,i),h(l,N,i),h(l,q,i),u(q,V);for(let n=0;n<$.length;n+=1)$[n]&&$[n].m(V,null);u(q,H),u(q,I);for(let n=0;n<p.length;n+=1)p[n]&&p[n].m(I,null);A=!0},p(l,[i]){(!A||i&1)&&d!==(d=l[0].name+"")&&Y(g,d),i&6&&(_=F(l[2]),$=x($,i,X,1,l,_,J,V,Te,pe,null,me)),i&6&&(E=F(l[2]),ee(),p=x(p,i,W,1,l,E,T,I,Ce,_e,null,de),te())},i(l){if(!A){for(let i=0;i<E.length;i+=1)U(p[i]);A=!0}},o(l){for(let i=0;i<p.length;i+=1)j(p[i]);A=!1},d(l){l&&(b(t),b(D),b(P),b(L),b(R),b(B),b(O),b(N),b(q));for(let i=0;i<$.length;i+=1)$[i].d();for(let i=0;i<p.length;i+=1)p[i].d()}}}function Re(s,t,e){let{collection:o}=t,f=204,c=[];const r=a=>e(1,f=a.code);return s.$$set=a=>{"collection"in a&&e(0,o=a.collection)},e(2,c=[{code:204,body:"null"},{code:400,body:`
                {
                  "status": 400,
                  "message": "An error occurred while validating the submitted data.",
                  "data": {
                    "token": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `}]),[o,f,c,r]}class Be extends le{constructor(t){super(),ne(this,t,Re,Pe,ie,{collection:0})}}function be(s,t,e){const o=s.slice();return o[4]=t[e],o}function he(s,t,e){const o=s.slice();return o[4]=t[e],o}function ve(s,t){let e,o=t[4].code+"",f,c,r,a;function d(){return t[3](t[4])}return{key:s,first:null,c(){e=m("button"),f=M(o),c=y(),v(e,"class","tab-item"),K(e,"active",t[1]===t[4].code),this.first=e},m(g,C){h(g,e,C),u(e,f),u(e,c),r||(a=oe(e,"click",d),r=!0)},p(g,C){t=g,C&4&&o!==(o=t[4].code+"")&&Y(f,o),C&6&&K(e,"active",t[1]===t[4].code)},d(g){g&&b(e),r=!1,a()}}}function ge(s,t){let e,o,f,c;return o=new qe({props:{content:t[4].body}}),{key:s,first:null,c(){e=m("div"),Q(o.$$.fragment),f=y(),v(e,"class","tab-item"),K(e,"active",t[1]===t[4].code),this.first=e},m(r,a){h(r,e,a),G(o,e,null),u(e,f),c=!0},p(r,a){t=r;const d={};a&4&&(d.content=t[4].body),o.$set(d),(!c||a&6)&&K(e,"active",t[1]===t[4].code)},i(r){c||(U(o.$$.fragment,r),c=!0)},o(r){j(o.$$.fragment,r),c=!1},d(r){r&&b(e),z(o)}}}function Oe(s){let t,e,o,f,c,r,a,d=s[0].name+"",g,C,D,P,L,R,B,O,N,q,V,$=[],J=new Map,H,I,p=[],T=new Map,A,_=F(s[2]);const X=l=>l[4].code;for(let l=0;l<_.length;l+=1){let i=he(s,_,l),n=X(i);J.set(n,$[l]=ve(n,i))}let E=F(s[2]);const W=l=>l[4].code;for(let l=0;l<E.length;l+=1){let i=be(s,E,l),n=W(i);T.set(n,p[l]=ge(n,i))}return{c(){t=m("div"),e=m("strong"),e.textContent="POST",o=y(),f=m("div"),c=m("p"),r=M("/api/collections/"),a=m("strong"),g=M(d),C=M("/request-verification"),D=y(),P=m("div"),P.textContent="Body Parameters",L=y(),R=m("table"),R.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>email</span></div></td> <td><span class="label">String</span></td> <td>The auth record email address to send the verification request (if exists).</td></tr></tbody>',B=y(),O=m("div"),O.textContent="Responses",N=y(),q=m("div"),V=m("div");for(let l=0;l<$.length;l+=1)$[l].c();H=y(),I=m("div");for(let l=0;l<p.length;l+=1)p[l].c();v(e,"class","label label-primary"),v(f,"class","content"),v(t,"class","alert alert-success"),v(P,"class","section-title"),v(R,"class","table-compact table-border m-b-base"),v(O,"class","section-title"),v(V,"class","tabs-header compact combined left"),v(I,"class","tabs-content"),v(q,"class","tabs")},m(l,i){h(l,t,i),u(t,e),u(t,o),u(t,f),u(f,c),u(c,r),u(c,a),u(a,g),u(c,C),h(l,D,i),h(l,P,i),h(l,L,i),h(l,R,i),h(l,B,i),h(l,O,i),h(l,N,i),h(l,q,i),u(q,V);for(let n=0;n<$.length;n+=1)$[n]&&$[n].m(V,null);u(q,H),u(q,I);for(let n=0;n<p.length;n+=1)p[n]&&p[n].m(I,null);A=!0},p(l,[i]){(!A||i&1)&&d!==(d=l[0].name+"")&&Y(g,d),i&6&&(_=F(l[2]),$=x($,i,X,1,l,_,J,V,Te,ve,null,he)),i&6&&(E=F(l[2]),ee(),p=x(p,i,W,1,l,E,T,I,Ce,ge,null,be),te())},i(l){if(!A){for(let i=0;i<E.length;i+=1)U(p[i]);A=!0}},o(l){for(let i=0;i<p.length;i+=1)j(p[i]);A=!1},d(l){l&&(b(t),b(D),b(P),b(L),b(R),b(B),b(O),b(N),b(q));for(let i=0;i<$.length;i+=1)$[i].d();for(let i=0;i<p.length;i+=1)p[i].d()}}}function Ee(s,t,e){let{collection:o}=t,f=204,c=[];const r=a=>e(1,f=a.code);return s.$$set=a=>{"collection"in a&&e(0,o=a.collection)},e(2,c=[{code:204,body:"null"},{code:400,body:`
                {
                  "status": 400,
                  "message": "An error occurred while validating the submitted data.",
                  "data": {
                    "email": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `}]),[o,f,c,r]}class Me extends le{constructor(t){super(),ne(this,t,Ee,Oe,ie,{collection:0})}}function ke(s,t,e){const o=s.slice();return o[5]=t[e],o[7]=e,o}function $e(s,t,e){const o=s.slice();return o[5]=t[e],o[7]=e,o}function we(s){let t,e,o,f,c;function r(){return s[4](s[7])}return{c(){t=m("button"),e=m("div"),e.textContent=`${s[5].title}`,o=y(),v(e,"class","txt"),v(t,"class","tab-item"),K(t,"active",s[1]==s[7])},m(a,d){h(a,t,d),u(t,e),u(t,o),f||(c=oe(t,"click",r),f=!0)},p(a,d){s=a,d&2&&K(t,"active",s[1]==s[7])},d(a){a&&b(t),f=!1,c()}}}function ye(s){let t,e,o,f;var c=s[5].component;function r(a,d){return{props:{collection:a[0]}}}return c&&(e=ue(c,r(s))),{c(){t=m("div"),e&&Q(e.$$.fragment),o=y(),v(t,"class","tab-item"),K(t,"active",s[1]==s[7])},m(a,d){h(a,t,d),e&&G(e,t,null),u(t,o),f=!0},p(a,d){if(c!==(c=a[5].component)){if(e){ee();const g=e;j(g.$$.fragment,1,0,()=>{z(g,1)}),te()}c?(e=ue(c,r(a)),Q(e.$$.fragment),U(e.$$.fragment,1),G(e,t,o)):e=null}else if(c){const g={};d&1&&(g.collection=a[0]),e.$set(g)}(!f||d&2)&&K(t,"active",a[1]==a[7])},i(a){f||(e&&U(e.$$.fragment,a),f=!0)},o(a){e&&j(e.$$.fragment,a),f=!1},d(a){a&&b(t),e&&z(e)}}}function Ne(s){var E,W,l,i;let t,e,o=s[0].name+"",f,c,r,a,d,g,C,D=s[0].name+"",P,L,R,B,O,N,q,V,$,J,H,I;B=new Ve({props:{js:`
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
    `}});let p=F(s[3]),T=[];for(let n=0;n<p.length;n+=1)T[n]=we($e(s,p,n));let A=F(s[3]),_=[];for(let n=0;n<A.length;n+=1)_[n]=ye(ke(s,A,n));const X=n=>j(_[n],1,1,()=>{_[n]=null});return{c(){t=m("h3"),e=M("Account verification ("),f=M(o),c=M(")"),r=y(),a=m("div"),d=m("p"),g=M("Sends "),C=m("strong"),P=M(D),L=M(" account verification request."),R=y(),Q(B.$$.fragment),O=y(),N=m("h6"),N.textContent="API details",q=y(),V=m("div"),$=m("div");for(let n=0;n<T.length;n+=1)T[n].c();J=y(),H=m("div");for(let n=0;n<_.length;n+=1)_[n].c();v(t,"class","m-b-sm"),v(a,"class","content txt-lg m-b-sm"),v(N,"class","m-b-xs"),v($,"class","tabs-header compact"),v(H,"class","tabs-content"),v(V,"class","tabs")},m(n,w){h(n,t,w),u(t,e),u(t,f),u(t,c),h(n,r,w),h(n,a,w),u(a,d),u(d,g),u(d,C),u(C,P),u(d,L),h(n,R,w),G(B,n,w),h(n,O,w),h(n,N,w),h(n,q,w),h(n,V,w),u(V,$);for(let S=0;S<T.length;S+=1)T[S]&&T[S].m($,null);u(V,J),u(V,H);for(let S=0;S<_.length;S+=1)_[S]&&_[S].m(H,null);I=!0},p(n,[w]){var se,ae,ce,re;(!I||w&1)&&o!==(o=n[0].name+"")&&Y(f,o),(!I||w&1)&&D!==(D=n[0].name+"")&&Y(P,D);const S={};if(w&5&&(S.js=`
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
    `),B.$set(S),w&10){p=F(n[3]);let k;for(k=0;k<p.length;k+=1){const Z=$e(n,p,k);T[k]?T[k].p(Z,w):(T[k]=we(Z),T[k].c(),T[k].m($,null))}for(;k<T.length;k+=1)T[k].d(1);T.length=p.length}if(w&11){A=F(n[3]);let k;for(k=0;k<A.length;k+=1){const Z=ke(n,A,k);_[k]?(_[k].p(Z,w),U(_[k],1)):(_[k]=ye(Z),_[k].c(),U(_[k],1),_[k].m(H,null))}for(ee(),k=A.length;k<_.length;k+=1)X(k);te()}},i(n){if(!I){U(B.$$.fragment,n);for(let w=0;w<A.length;w+=1)U(_[w]);I=!0}},o(n){j(B.$$.fragment,n),_=_.filter(Boolean);for(let w=0;w<_.length;w+=1)j(_[w]);I=!1},d(n){n&&(b(t),b(r),b(a),b(R),b(O),b(N),b(q),b(V)),z(B,n),fe(T,n),fe(_,n)}}}function Se(s,t,e){let o,{collection:f}=t;const c=[{title:"Request verification",component:Me},{title:"Confirm verification",component:Be}];let r=0;const a=d=>e(1,r=d);return s.$$set=d=>{"collection"in d&&e(0,f=d.collection)},e(2,o=Ae.getApiExampleUrl(Ie.baseURL)),[f,r,o,c,a]}class Fe extends le{constructor(t){super(),ne(this,t,Se,Ne,ie,{collection:0})}}export{Fe as default};
