import{S as se,i as oe,s as ie,X as K,h as g,t as X,a as V,I as F,Z as le,_ as Re,C as ne,$ as Se,D as ae,l as v,n as u,u as p,v as y,A as U,w as b,k as Y,o as ce,W as Oe,d as x,m as ee,c as te,V as Me,Y as _e,J as Be,p as De,a0 as be}from"./index-pGELYd11.js";function ge(n,e,t){const l=n.slice();return l[4]=e[t],l}function ve(n,e,t){const l=n.slice();return l[4]=e[t],l}function ke(n,e){let t,l=e[4].code+"",d,i,r,a;function m(){return e[3](e[4])}return{key:n,first:null,c(){t=p("button"),d=U(l),i=y(),b(t,"class","tab-item"),Y(t,"active",e[1]===e[4].code),this.first=t},m(k,q){v(k,t,q),u(t,d),u(t,i),r||(a=ce(t,"click",m),r=!0)},p(k,q){e=k,q&4&&l!==(l=e[4].code+"")&&F(d,l),q&6&&Y(t,"active",e[1]===e[4].code)},d(k){k&&g(t),r=!1,a()}}}function $e(n,e){let t,l,d,i;return l=new Oe({props:{content:e[4].body}}),{key:n,first:null,c(){t=p("div"),te(l.$$.fragment),d=y(),b(t,"class","tab-item"),Y(t,"active",e[1]===e[4].code),this.first=t},m(r,a){v(r,t,a),ee(l,t,null),u(t,d),i=!0},p(r,a){e=r;const m={};a&4&&(m.content=e[4].body),l.$set(m),(!i||a&6)&&Y(t,"active",e[1]===e[4].code)},i(r){i||(V(l.$$.fragment,r),i=!0)},o(r){X(l.$$.fragment,r),i=!1},d(r){r&&g(t),x(l)}}}function Ne(n){let e,t,l,d,i,r,a,m=n[0].name+"",k,q,G,H,J,L,z,B,D,S,N,A=[],O=new Map,P,j,T=[],W=new Map,w,E=K(n[2]);const M=c=>c[4].code;for(let c=0;c<E.length;c+=1){let f=ve(n,E,c),s=M(f);O.set(s,A[c]=ke(s,f))}let _=K(n[2]);const Z=c=>c[4].code;for(let c=0;c<_.length;c+=1){let f=ge(n,_,c),s=Z(f);W.set(s,T[c]=$e(s,f))}return{c(){e=p("div"),t=p("strong"),t.textContent="POST",l=y(),d=p("div"),i=p("p"),r=U("/api/collections/"),a=p("strong"),k=U(m),q=U("/confirm-email-change"),G=y(),H=p("div"),H.textContent="Body Parameters",J=y(),L=p("table"),L.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>token</span></div></td> <td><span class="label">String</span></td> <td>The token from the change email request email.</td></tr> <tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>The account password to confirm the email change.</td></tr></tbody>',z=y(),B=p("div"),B.textContent="Responses",D=y(),S=p("div"),N=p("div");for(let c=0;c<A.length;c+=1)A[c].c();P=y(),j=p("div");for(let c=0;c<T.length;c+=1)T[c].c();b(t,"class","label label-primary"),b(d,"class","content"),b(e,"class","alert alert-success"),b(H,"class","section-title"),b(L,"class","table-compact table-border m-b-base"),b(B,"class","section-title"),b(N,"class","tabs-header compact combined left"),b(j,"class","tabs-content"),b(S,"class","tabs")},m(c,f){v(c,e,f),u(e,t),u(e,l),u(e,d),u(d,i),u(i,r),u(i,a),u(a,k),u(i,q),v(c,G,f),v(c,H,f),v(c,J,f),v(c,L,f),v(c,z,f),v(c,B,f),v(c,D,f),v(c,S,f),u(S,N);for(let s=0;s<A.length;s+=1)A[s]&&A[s].m(N,null);u(S,P),u(S,j);for(let s=0;s<T.length;s+=1)T[s]&&T[s].m(j,null);w=!0},p(c,[f]){(!w||f&1)&&m!==(m=c[0].name+"")&&F(k,m),f&6&&(E=K(c[2]),A=le(A,f,M,1,c,E,O,N,Re,ke,null,ve)),f&6&&(_=K(c[2]),ne(),T=le(T,f,Z,1,c,_,W,j,Se,$e,null,ge),ae())},i(c){if(!w){for(let f=0;f<_.length;f+=1)V(T[f]);w=!0}},o(c){for(let f=0;f<T.length;f+=1)X(T[f]);w=!1},d(c){c&&(g(e),g(G),g(H),g(J),g(L),g(z),g(B),g(D),g(S));for(let f=0;f<A.length;f+=1)A[f].d();for(let f=0;f<T.length;f+=1)T[f].d()}}}function We(n,e,t){let{collection:l}=e,d=204,i=[];const r=a=>t(1,d=a.code);return n.$$set=a=>{"collection"in a&&t(0,l=a.collection)},t(2,i=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[l,d,i,r]}class He extends se{constructor(e){super(),oe(this,e,We,Ne,ie,{collection:0})}}function we(n,e,t){const l=n.slice();return l[4]=e[t],l}function Ce(n,e,t){const l=n.slice();return l[4]=e[t],l}function ye(n,e){let t,l=e[4].code+"",d,i,r,a;function m(){return e[3](e[4])}return{key:n,first:null,c(){t=p("button"),d=U(l),i=y(),b(t,"class","tab-item"),Y(t,"active",e[1]===e[4].code),this.first=t},m(k,q){v(k,t,q),u(t,d),u(t,i),r||(a=ce(t,"click",m),r=!0)},p(k,q){e=k,q&4&&l!==(l=e[4].code+"")&&F(d,l),q&6&&Y(t,"active",e[1]===e[4].code)},d(k){k&&g(t),r=!1,a()}}}function Ee(n,e){let t,l,d,i;return l=new Oe({props:{content:e[4].body}}),{key:n,first:null,c(){t=p("div"),te(l.$$.fragment),d=y(),b(t,"class","tab-item"),Y(t,"active",e[1]===e[4].code),this.first=t},m(r,a){v(r,t,a),ee(l,t,null),u(t,d),i=!0},p(r,a){e=r;const m={};a&4&&(m.content=e[4].body),l.$set(m),(!i||a&6)&&Y(t,"active",e[1]===e[4].code)},i(r){i||(V(l.$$.fragment,r),i=!0)},o(r){X(l.$$.fragment,r),i=!1},d(r){r&&g(t),x(l)}}}function Le(n){let e,t,l,d,i,r,a,m=n[0].name+"",k,q,G,H,J,L,z,B,D,S,N,A,O,P=[],j=new Map,T,W,w=[],E=new Map,M,_=K(n[2]);const Z=s=>s[4].code;for(let s=0;s<_.length;s+=1){let h=Ce(n,_,s),R=Z(h);j.set(R,P[s]=ye(R,h))}let c=K(n[2]);const f=s=>s[4].code;for(let s=0;s<c.length;s+=1){let h=we(n,c,s),R=f(h);E.set(R,w[s]=Ee(R,h))}return{c(){e=p("div"),t=p("strong"),t.textContent="POST",l=y(),d=p("div"),i=p("p"),r=U("/api/collections/"),a=p("strong"),k=U(m),q=U("/request-email-change"),G=y(),H=p("p"),H.innerHTML="Requires <code>Authorization:TOKEN</code>",J=y(),L=p("div"),L.textContent="Body Parameters",z=y(),B=p("table"),B.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>newEmail</span></div></td> <td><span class="label">String</span></td> <td>The new email address to send the change email request.</td></tr></tbody>',D=y(),S=p("div"),S.textContent="Responses",N=y(),A=p("div"),O=p("div");for(let s=0;s<P.length;s+=1)P[s].c();T=y(),W=p("div");for(let s=0;s<w.length;s+=1)w[s].c();b(t,"class","label label-primary"),b(d,"class","content"),b(H,"class","txt-hint txt-sm txt-right"),b(e,"class","alert alert-success"),b(L,"class","section-title"),b(B,"class","table-compact table-border m-b-base"),b(S,"class","section-title"),b(O,"class","tabs-header compact combined left"),b(W,"class","tabs-content"),b(A,"class","tabs")},m(s,h){v(s,e,h),u(e,t),u(e,l),u(e,d),u(d,i),u(i,r),u(i,a),u(a,k),u(i,q),u(e,G),u(e,H),v(s,J,h),v(s,L,h),v(s,z,h),v(s,B,h),v(s,D,h),v(s,S,h),v(s,N,h),v(s,A,h),u(A,O);for(let R=0;R<P.length;R+=1)P[R]&&P[R].m(O,null);u(A,T),u(A,W);for(let R=0;R<w.length;R+=1)w[R]&&w[R].m(W,null);M=!0},p(s,[h]){(!M||h&1)&&m!==(m=s[0].name+"")&&F(k,m),h&6&&(_=K(s[2]),P=le(P,h,Z,1,s,_,j,O,Re,ye,null,Ce)),h&6&&(c=K(s[2]),ne(),w=le(w,h,f,1,s,c,E,W,Se,Ee,null,we),ae())},i(s){if(!M){for(let h=0;h<c.length;h+=1)V(w[h]);M=!0}},o(s){for(let h=0;h<w.length;h+=1)X(w[h]);M=!1},d(s){s&&(g(e),g(J),g(L),g(z),g(B),g(D),g(S),g(N),g(A));for(let h=0;h<P.length;h+=1)P[h].d();for(let h=0;h<w.length;h+=1)w[h].d()}}}function Ue(n,e,t){let{collection:l}=e,d=204,i=[];const r=a=>t(1,d=a.code);return n.$$set=a=>{"collection"in a&&t(0,l=a.collection)},t(2,i=[{code:204,body:"null"},{code:400,body:`
                {
                  "status": 400,
                  "message": "An error occurred while validating the submitted data.",
                  "data": {
                    "newEmail": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `},{code:401,body:`
                {
                  "status": 401,
                  "message": "The request requires valid record authorization token to be set.",
                  "data": {}
                }
            `},{code:403,body:`
                {
                  "status": 403,
                  "message": "The authorized record model is not allowed to perform this action.",
                  "data": {}
                }
            `}]),[l,d,i,r]}class Ie extends se{constructor(e){super(),oe(this,e,Ue,Le,ie,{collection:0})}}function Ae(n,e,t){const l=n.slice();return l[5]=e[t],l[7]=t,l}function Te(n,e,t){const l=n.slice();return l[5]=e[t],l[7]=t,l}function qe(n){let e,t,l,d,i;function r(){return n[4](n[7])}return{c(){e=p("button"),t=p("div"),t.textContent=`${n[5].title}`,l=y(),b(t,"class","txt"),b(e,"class","tab-item"),Y(e,"active",n[1]==n[7])},m(a,m){v(a,e,m),u(e,t),u(e,l),d||(i=ce(e,"click",r),d=!0)},p(a,m){n=a,m&2&&Y(e,"active",n[1]==n[7])},d(a){a&&g(e),d=!1,i()}}}function Pe(n){let e,t,l,d;var i=n[5].component;function r(a,m){return{props:{collection:a[0]}}}return i&&(t=be(i,r(n))),{c(){e=p("div"),t&&te(t.$$.fragment),l=y(),b(e,"class","tab-item"),Y(e,"active",n[1]==n[7])},m(a,m){v(a,e,m),t&&ee(t,e,null),u(e,l),d=!0},p(a,m){if(i!==(i=a[5].component)){if(t){ne();const k=t;X(k.$$.fragment,1,0,()=>{x(k,1)}),ae()}i?(t=be(i,r(a)),te(t.$$.fragment),V(t.$$.fragment,1),ee(t,e,l)):t=null}else if(i){const k={};m&1&&(k.collection=a[0]),t.$set(k)}(!d||m&2)&&Y(e,"active",a[1]==a[7])},i(a){d||(t&&V(t.$$.fragment,a),d=!0)},o(a){t&&X(t.$$.fragment,a),d=!1},d(a){a&&g(e),t&&x(t)}}}function Ke(n){var c,f,s,h,R,re;let e,t,l=n[0].name+"",d,i,r,a,m,k,q,G=n[0].name+"",H,J,L,z,B,D,S,N,A,O,P,j,T,W;D=new Me({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${n[2]}');

        ...

        await pb.collection('${(c=n[0])==null?void 0:c.name}').authWithPassword('test@example.com', '1234567890');

        await pb.collection('${(f=n[0])==null?void 0:f.name}').requestEmailChange('new@example.com');

        // ---
        // (optional) in your custom confirmation page:
        // ---

        // note: after this call all previously issued auth tokens are invalidated
        await pb.collection('${(s=n[0])==null?void 0:s.name}').confirmEmailChange(
            'EMAIL_CHANGE_TOKEN',
            'YOUR_PASSWORD',
        );
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${n[2]}');

        ...

        await pb.collection('${(h=n[0])==null?void 0:h.name}').authWithPassword('test@example.com', '1234567890');

        await pb.collection('${(R=n[0])==null?void 0:R.name}').requestEmailChange('new@example.com');

        ...

        // ---
        // (optional) in your custom confirmation page:
        // ---

        // note: after this call all previously issued auth tokens are invalidated
        await pb.collection('${(re=n[0])==null?void 0:re.name}').confirmEmailChange(
          'EMAIL_CHANGE_TOKEN',
          'YOUR_PASSWORD',
        );
    `}});let w=K(n[3]),E=[];for(let o=0;o<w.length;o+=1)E[o]=qe(Te(n,w,o));let M=K(n[3]),_=[];for(let o=0;o<M.length;o+=1)_[o]=Pe(Ae(n,M,o));const Z=o=>X(_[o],1,1,()=>{_[o]=null});return{c(){e=p("h3"),t=U("Email change ("),d=U(l),i=U(")"),r=y(),a=p("div"),m=p("p"),k=U("Sends "),q=p("strong"),H=U(G),J=U(" email change request."),L=y(),z=p("p"),z.textContent=`On successful email change all previously issued auth tokens for the specific record will be
        automatically invalidated.`,B=y(),te(D.$$.fragment),S=y(),N=p("h6"),N.textContent="API details",A=y(),O=p("div"),P=p("div");for(let o=0;o<E.length;o+=1)E[o].c();j=y(),T=p("div");for(let o=0;o<_.length;o+=1)_[o].c();b(e,"class","m-b-sm"),b(a,"class","content txt-lg m-b-sm"),b(N,"class","m-b-xs"),b(P,"class","tabs-header compact"),b(T,"class","tabs-content"),b(O,"class","tabs")},m(o,C){v(o,e,C),u(e,t),u(e,d),u(e,i),v(o,r,C),v(o,a,C),u(a,m),u(m,k),u(m,q),u(q,H),u(m,J),u(a,L),u(a,z),v(o,B,C),ee(D,o,C),v(o,S,C),v(o,N,C),v(o,A,C),v(o,O,C),u(O,P);for(let I=0;I<E.length;I+=1)E[I]&&E[I].m(P,null);u(O,j),u(O,T);for(let I=0;I<_.length;I+=1)_[I]&&_[I].m(T,null);W=!0},p(o,[C]){var de,ue,fe,me,he,pe;(!W||C&1)&&l!==(l=o[0].name+"")&&F(d,l),(!W||C&1)&&G!==(G=o[0].name+"")&&F(H,G);const I={};if(C&5&&(I.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${o[2]}');

        ...

        await pb.collection('${(de=o[0])==null?void 0:de.name}').authWithPassword('test@example.com', '1234567890');

        await pb.collection('${(ue=o[0])==null?void 0:ue.name}').requestEmailChange('new@example.com');

        // ---
        // (optional) in your custom confirmation page:
        // ---

        // note: after this call all previously issued auth tokens are invalidated
        await pb.collection('${(fe=o[0])==null?void 0:fe.name}').confirmEmailChange(
            'EMAIL_CHANGE_TOKEN',
            'YOUR_PASSWORD',
        );
    `),C&5&&(I.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${o[2]}');

        ...

        await pb.collection('${(me=o[0])==null?void 0:me.name}').authWithPassword('test@example.com', '1234567890');

        await pb.collection('${(he=o[0])==null?void 0:he.name}').requestEmailChange('new@example.com');

        ...

        // ---
        // (optional) in your custom confirmation page:
        // ---

        // note: after this call all previously issued auth tokens are invalidated
        await pb.collection('${(pe=o[0])==null?void 0:pe.name}').confirmEmailChange(
          'EMAIL_CHANGE_TOKEN',
          'YOUR_PASSWORD',
        );
    `),D.$set(I),C&10){w=K(o[3]);let $;for($=0;$<w.length;$+=1){const Q=Te(o,w,$);E[$]?E[$].p(Q,C):(E[$]=qe(Q),E[$].c(),E[$].m(P,null))}for(;$<E.length;$+=1)E[$].d(1);E.length=w.length}if(C&11){M=K(o[3]);let $;for($=0;$<M.length;$+=1){const Q=Ae(o,M,$);_[$]?(_[$].p(Q,C),V(_[$],1)):(_[$]=Pe(Q),_[$].c(),V(_[$],1),_[$].m(T,null))}for(ne(),$=M.length;$<_.length;$+=1)Z($);ae()}},i(o){if(!W){V(D.$$.fragment,o);for(let C=0;C<M.length;C+=1)V(_[C]);W=!0}},o(o){X(D.$$.fragment,o),_=_.filter(Boolean);for(let C=0;C<_.length;C+=1)X(_[C]);W=!1},d(o){o&&(g(e),g(r),g(a),g(B),g(S),g(N),g(A),g(O)),x(D,o),_e(E,o),_e(_,o)}}}function Ye(n,e,t){let l,{collection:d}=e;const i=[{title:"Request email change",component:Ie},{title:"Confirm email change",component:He}];let r=0;const a=m=>t(1,r=m);return n.$$set=m=>{"collection"in m&&t(0,d=m.collection)},t(2,l=Be.getApiExampleUrl(De.baseURL)),[d,r,l,i,a]}class ze extends se{constructor(e){super(),oe(this,e,Ye,Ke,ie,{collection:0})}}export{ze as default};
