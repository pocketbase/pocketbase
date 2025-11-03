import{S as se,i as ne,s as ae,O as H,h as b,t as z,a as Q,I as V,ad as ee,ae as ye,C as te,af as Te,D as le,l as v,n as u,u as p,v as S,A as D,w as k,k as L,o as oe,ac as Ee,d as Y,m as Z,c as x,ab as Ce,Q as fe,J as Oe,p as qe,a4 as pe}from"./index-B4ZsHsKR.js";function me(a,t,e){const n=a.slice();return n[4]=t[e],n}function _e(a,t,e){const n=a.slice();return n[4]=t[e],n}function he(a,t){let e,n=t[4].code+"",d,c,r,o;function f(){return t[3](t[4])}return{key:a,first:null,c(){e=p("button"),d=D(n),c=S(),k(e,"class","tab-item"),L(e,"active",t[1]===t[4].code),this.first=e},m(g,y){v(g,e,y),u(e,d),u(e,c),r||(o=oe(e,"click",f),r=!0)},p(g,y){t=g,y&4&&n!==(n=t[4].code+"")&&V(d,n),y&6&&L(e,"active",t[1]===t[4].code)},d(g){g&&b(e),r=!1,o()}}}function be(a,t){let e,n,d,c;return n=new Ee({props:{content:t[4].body}}),{key:a,first:null,c(){e=p("div"),x(n.$$.fragment),d=S(),k(e,"class","tab-item"),L(e,"active",t[1]===t[4].code),this.first=e},m(r,o){v(r,e,o),Z(n,e,null),u(e,d),c=!0},p(r,o){t=r;const f={};o&4&&(f.content=t[4].body),n.$set(f),(!c||o&6)&&L(e,"active",t[1]===t[4].code)},i(r){c||(Q(n.$$.fragment,r),c=!0)},o(r){z(n.$$.fragment,r),c=!1},d(r){r&&b(e),Y(n)}}}function Ae(a){let t,e,n,d,c,r,o,f=a[0].name+"",g,y,F,O,J,N,U,q,A,T,C,R=[],M=new Map,j,W,h=[],K=new Map,E,P=H(a[2]);const B=l=>l[4].code;for(let l=0;l<P.length;l+=1){let s=_e(a,P,l),_=B(s);M.set(_,R[l]=he(_,s))}let m=H(a[2]);const G=l=>l[4].code;for(let l=0;l<m.length;l+=1){let s=me(a,m,l),_=G(s);K.set(_,h[l]=be(_,s))}return{c(){t=p("div"),e=p("strong"),e.textContent="POST",n=S(),d=p("div"),c=p("p"),r=D("/api/collections/"),o=p("strong"),g=D(f),y=D("/confirm-password-reset"),F=S(),O=p("div"),O.textContent="Body Parameters",J=S(),N=p("table"),N.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>token</span></div></td> <td><span class="label">String</span></td> <td>The token from the password reset request email.</td></tr> <tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>The new password to set.</td></tr> <tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>passwordConfirm</span></div></td> <td><span class="label">String</span></td> <td>The new password confirmation.</td></tr></tbody>',U=S(),q=p("div"),q.textContent="Responses",A=S(),T=p("div"),C=p("div");for(let l=0;l<R.length;l+=1)R[l].c();j=S(),W=p("div");for(let l=0;l<h.length;l+=1)h[l].c();k(e,"class","label label-primary"),k(d,"class","content"),k(t,"class","alert alert-success"),k(O,"class","section-title"),k(N,"class","table-compact table-border m-b-base"),k(q,"class","section-title"),k(C,"class","tabs-header compact combined left"),k(W,"class","tabs-content"),k(T,"class","tabs")},m(l,s){v(l,t,s),u(t,e),u(t,n),u(t,d),u(d,c),u(c,r),u(c,o),u(o,g),u(c,y),v(l,F,s),v(l,O,s),v(l,J,s),v(l,N,s),v(l,U,s),v(l,q,s),v(l,A,s),v(l,T,s),u(T,C);for(let _=0;_<R.length;_+=1)R[_]&&R[_].m(C,null);u(T,j),u(T,W);for(let _=0;_<h.length;_+=1)h[_]&&h[_].m(W,null);E=!0},p(l,[s]){(!E||s&1)&&f!==(f=l[0].name+"")&&V(g,f),s&6&&(P=H(l[2]),R=ee(R,s,B,1,l,P,M,C,ye,he,null,_e)),s&6&&(m=H(l[2]),te(),h=ee(h,s,G,1,l,m,K,W,Te,be,null,me),le())},i(l){if(!E){for(let s=0;s<m.length;s+=1)Q(h[s]);E=!0}},o(l){for(let s=0;s<h.length;s+=1)z(h[s]);E=!1},d(l){l&&(b(t),b(F),b(O),b(J),b(N),b(U),b(q),b(A),b(T));for(let s=0;s<R.length;s+=1)R[s].d();for(let s=0;s<h.length;s+=1)h[s].d()}}}function Ne(a,t,e){let{collection:n}=t,d=204,c=[];const r=o=>e(1,d=o.code);return a.$$set=o=>{"collection"in o&&e(0,n=o.collection)},e(2,c=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[n,d,c,r]}class We extends se{constructor(t){super(),ne(this,t,Ne,Ae,ae,{collection:0})}}function ve(a,t,e){const n=a.slice();return n[4]=t[e],n}function ke(a,t,e){const n=a.slice();return n[4]=t[e],n}function ge(a,t){let e,n=t[4].code+"",d,c,r,o;function f(){return t[3](t[4])}return{key:a,first:null,c(){e=p("button"),d=D(n),c=S(),k(e,"class","tab-item"),L(e,"active",t[1]===t[4].code),this.first=e},m(g,y){v(g,e,y),u(e,d),u(e,c),r||(o=oe(e,"click",f),r=!0)},p(g,y){t=g,y&4&&n!==(n=t[4].code+"")&&V(d,n),y&6&&L(e,"active",t[1]===t[4].code)},d(g){g&&b(e),r=!1,o()}}}function we(a,t){let e,n,d,c;return n=new Ee({props:{content:t[4].body}}),{key:a,first:null,c(){e=p("div"),x(n.$$.fragment),d=S(),k(e,"class","tab-item"),L(e,"active",t[1]===t[4].code),this.first=e},m(r,o){v(r,e,o),Z(n,e,null),u(e,d),c=!0},p(r,o){t=r;const f={};o&4&&(f.content=t[4].body),n.$set(f),(!c||o&6)&&L(e,"active",t[1]===t[4].code)},i(r){c||(Q(n.$$.fragment,r),c=!0)},o(r){z(n.$$.fragment,r),c=!1},d(r){r&&b(e),Y(n)}}}function De(a){let t,e,n,d,c,r,o,f=a[0].name+"",g,y,F,O,J,N,U,q,A,T,C,R=[],M=new Map,j,W,h=[],K=new Map,E,P=H(a[2]);const B=l=>l[4].code;for(let l=0;l<P.length;l+=1){let s=ke(a,P,l),_=B(s);M.set(_,R[l]=ge(_,s))}let m=H(a[2]);const G=l=>l[4].code;for(let l=0;l<m.length;l+=1){let s=ve(a,m,l),_=G(s);K.set(_,h[l]=we(_,s))}return{c(){t=p("div"),e=p("strong"),e.textContent="POST",n=S(),d=p("div"),c=p("p"),r=D("/api/collections/"),o=p("strong"),g=D(f),y=D("/request-password-reset"),F=S(),O=p("div"),O.textContent="Body Parameters",J=S(),N=p("table"),N.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>email</span></div></td> <td><span class="label">String</span></td> <td>The auth record email address to send the password reset request (if exists).</td></tr></tbody>',U=S(),q=p("div"),q.textContent="Responses",A=S(),T=p("div"),C=p("div");for(let l=0;l<R.length;l+=1)R[l].c();j=S(),W=p("div");for(let l=0;l<h.length;l+=1)h[l].c();k(e,"class","label label-primary"),k(d,"class","content"),k(t,"class","alert alert-success"),k(O,"class","section-title"),k(N,"class","table-compact table-border m-b-base"),k(q,"class","section-title"),k(C,"class","tabs-header compact combined left"),k(W,"class","tabs-content"),k(T,"class","tabs")},m(l,s){v(l,t,s),u(t,e),u(t,n),u(t,d),u(d,c),u(c,r),u(c,o),u(o,g),u(c,y),v(l,F,s),v(l,O,s),v(l,J,s),v(l,N,s),v(l,U,s),v(l,q,s),v(l,A,s),v(l,T,s),u(T,C);for(let _=0;_<R.length;_+=1)R[_]&&R[_].m(C,null);u(T,j),u(T,W);for(let _=0;_<h.length;_+=1)h[_]&&h[_].m(W,null);E=!0},p(l,[s]){(!E||s&1)&&f!==(f=l[0].name+"")&&V(g,f),s&6&&(P=H(l[2]),R=ee(R,s,B,1,l,P,M,C,ye,ge,null,ke)),s&6&&(m=H(l[2]),te(),h=ee(h,s,G,1,l,m,K,W,Te,we,null,ve),le())},i(l){if(!E){for(let s=0;s<m.length;s+=1)Q(h[s]);E=!0}},o(l){for(let s=0;s<h.length;s+=1)z(h[s]);E=!1},d(l){l&&(b(t),b(F),b(O),b(J),b(N),b(U),b(q),b(A),b(T));for(let s=0;s<R.length;s+=1)R[s].d();for(let s=0;s<h.length;s+=1)h[s].d()}}}function Me(a,t,e){let{collection:n}=t,d=204,c=[];const r=o=>e(1,d=o.code);return a.$$set=o=>{"collection"in o&&e(0,n=o.collection)},e(2,c=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[n,d,c,r]}class Be extends se{constructor(t){super(),ne(this,t,Me,De,ae,{collection:0})}}function $e(a,t,e){const n=a.slice();return n[5]=t[e],n[7]=e,n}function Re(a,t,e){const n=a.slice();return n[5]=t[e],n[7]=e,n}function Pe(a){let t,e,n,d,c;function r(){return a[4](a[7])}return{c(){t=p("button"),e=p("div"),e.textContent=`${a[5].title}`,n=S(),k(e,"class","txt"),k(t,"class","tab-item"),L(t,"active",a[1]==a[7])},m(o,f){v(o,t,f),u(t,e),u(t,n),d||(c=oe(t,"click",r),d=!0)},p(o,f){a=o,f&2&&L(t,"active",a[1]==a[7])},d(o){o&&b(t),d=!1,c()}}}function Se(a){let t,e,n,d;var c=a[5].component;function r(o,f){return{props:{collection:o[0]}}}return c&&(e=pe(c,r(a))),{c(){t=p("div"),e&&x(e.$$.fragment),n=S(),k(t,"class","tab-item"),L(t,"active",a[1]==a[7])},m(o,f){v(o,t,f),e&&Z(e,t,null),u(t,n),d=!0},p(o,f){if(c!==(c=o[5].component)){if(e){te();const g=e;z(g.$$.fragment,1,0,()=>{Y(g,1)}),le()}c?(e=pe(c,r(o)),x(e.$$.fragment),Q(e.$$.fragment,1),Z(e,t,n)):e=null}else if(c){const g={};f&1&&(g.collection=o[0]),e.$set(g)}(!d||f&2)&&L(t,"active",o[1]==o[7])},i(o){d||(e&&Q(e.$$.fragment,o),d=!0)},o(o){e&&z(e.$$.fragment,o),d=!1},d(o){o&&b(t),e&&Y(e)}}}function Ie(a){var l,s,_,ie;let t,e,n=a[0].name+"",d,c,r,o,f,g,y,F=a[0].name+"",O,J,N,U,q,A,T,C,R,M,j,W,h,K;A=new Ce({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${a[2]}');

        ...

        await pb.collection('${(l=a[0])==null?void 0:l.name}').requestPasswordReset('test@example.com');

        // ---
        // (optional) in your custom confirmation page:
        // ---

        // note: after this call all previously issued auth tokens are invalidated
        await pb.collection('${(s=a[0])==null?void 0:s.name}').confirmPasswordReset(
            'RESET_TOKEN',
            'NEW_PASSWORD',
            'NEW_PASSWORD_CONFIRM',
        );
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${a[2]}');

        ...

        await pb.collection('${(_=a[0])==null?void 0:_.name}').requestPasswordReset('test@example.com');

        // ---
        // (optional) in your custom confirmation page:
        // ---

        // note: after this call all previously issued auth tokens are invalidated
        await pb.collection('${(ie=a[0])==null?void 0:ie.name}').confirmPasswordReset(
          'RESET_TOKEN',
          'NEW_PASSWORD',
          'NEW_PASSWORD_CONFIRM',
        );
    `}});let E=H(a[3]),P=[];for(let i=0;i<E.length;i+=1)P[i]=Pe(Re(a,E,i));let B=H(a[3]),m=[];for(let i=0;i<B.length;i+=1)m[i]=Se($e(a,B,i));const G=i=>z(m[i],1,1,()=>{m[i]=null});return{c(){t=p("h3"),e=D("Password reset ("),d=D(n),c=D(")"),r=S(),o=p("div"),f=p("p"),g=D("Sends "),y=p("strong"),O=D(F),J=D(" password reset email request."),N=S(),U=p("p"),U.textContent=`On successful password reset all previously issued auth tokens for the specific record will be
        automatically invalidated.`,q=S(),x(A.$$.fragment),T=S(),C=p("h6"),C.textContent="API details",R=S(),M=p("div"),j=p("div");for(let i=0;i<P.length;i+=1)P[i].c();W=S(),h=p("div");for(let i=0;i<m.length;i+=1)m[i].c();k(t,"class","m-b-sm"),k(o,"class","content txt-lg m-b-sm"),k(C,"class","m-b-xs"),k(j,"class","tabs-header compact"),k(h,"class","tabs-content"),k(M,"class","tabs")},m(i,$){v(i,t,$),u(t,e),u(t,d),u(t,c),v(i,r,$),v(i,o,$),u(o,f),u(f,g),u(f,y),u(y,O),u(f,J),u(o,N),u(o,U),v(i,q,$),Z(A,i,$),v(i,T,$),v(i,C,$),v(i,R,$),v(i,M,$),u(M,j);for(let I=0;I<P.length;I+=1)P[I]&&P[I].m(j,null);u(M,W),u(M,h);for(let I=0;I<m.length;I+=1)m[I]&&m[I].m(h,null);K=!0},p(i,[$]){var ce,re,de,ue;(!K||$&1)&&n!==(n=i[0].name+"")&&V(d,n),(!K||$&1)&&F!==(F=i[0].name+"")&&V(O,F);const I={};if($&5&&(I.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${i[2]}');

        ...

        await pb.collection('${(ce=i[0])==null?void 0:ce.name}').requestPasswordReset('test@example.com');

        // ---
        // (optional) in your custom confirmation page:
        // ---

        // note: after this call all previously issued auth tokens are invalidated
        await pb.collection('${(re=i[0])==null?void 0:re.name}').confirmPasswordReset(
            'RESET_TOKEN',
            'NEW_PASSWORD',
            'NEW_PASSWORD_CONFIRM',
        );
    `),$&5&&(I.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${i[2]}');

        ...

        await pb.collection('${(de=i[0])==null?void 0:de.name}').requestPasswordReset('test@example.com');

        // ---
        // (optional) in your custom confirmation page:
        // ---

        // note: after this call all previously issued auth tokens are invalidated
        await pb.collection('${(ue=i[0])==null?void 0:ue.name}').confirmPasswordReset(
          'RESET_TOKEN',
          'NEW_PASSWORD',
          'NEW_PASSWORD_CONFIRM',
        );
    `),A.$set(I),$&10){E=H(i[3]);let w;for(w=0;w<E.length;w+=1){const X=Re(i,E,w);P[w]?P[w].p(X,$):(P[w]=Pe(X),P[w].c(),P[w].m(j,null))}for(;w<P.length;w+=1)P[w].d(1);P.length=E.length}if($&11){B=H(i[3]);let w;for(w=0;w<B.length;w+=1){const X=$e(i,B,w);m[w]?(m[w].p(X,$),Q(m[w],1)):(m[w]=Se(X),m[w].c(),Q(m[w],1),m[w].m(h,null))}for(te(),w=B.length;w<m.length;w+=1)G(w);le()}},i(i){if(!K){Q(A.$$.fragment,i);for(let $=0;$<B.length;$+=1)Q(m[$]);K=!0}},o(i){z(A.$$.fragment,i),m=m.filter(Boolean);for(let $=0;$<m.length;$+=1)z(m[$]);K=!1},d(i){i&&(b(t),b(r),b(o),b(q),b(T),b(C),b(R),b(M)),Y(A,i),fe(P,i),fe(m,i)}}}function Fe(a,t,e){let n,{collection:d}=t;const c=[{title:"Request password reset",component:Be},{title:"Confirm password reset",component:We}];let r=0;const o=f=>e(1,r=f);return a.$$set=f=>{"collection"in f&&e(0,d=f.collection)},e(2,n=Oe.getApiExampleUrl(qe.baseURL)),[d,r,n,c,o]}class He extends se{constructor(t){super(),ne(this,t,Fe,Ie,ae,{collection:0})}}export{He as default};
