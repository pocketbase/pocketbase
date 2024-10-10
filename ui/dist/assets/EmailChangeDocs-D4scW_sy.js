import{S as se,i as oe,s as ie,T as K,e as p,b as y,w as U,f as b,g,h as u,x as J,U as le,V as Re,k as ne,W as Se,n as ae,t as Q,a as V,o as v,r as Y,u as ce,R as Oe,c as x,m as ee,d as te,Q as Me,X as _e,C as Be,p as De,Y as be}from"./index-C36u8XUk.js";function ge(n,e,t){const l=n.slice();return l[4]=e[t],l}function ve(n,e,t){const l=n.slice();return l[4]=e[t],l}function ke(n,e){let t,l=e[4].code+"",d,i,r,a;function m(){return e[3](e[4])}return{key:n,first:null,c(){t=p("button"),d=U(l),i=y(),b(t,"class","tab-item"),Y(t,"active",e[1]===e[4].code),this.first=t},m(k,P){g(k,t,P),u(t,d),u(t,i),r||(a=ce(t,"click",m),r=!0)},p(k,P){e=k,P&4&&l!==(l=e[4].code+"")&&J(d,l),P&6&&Y(t,"active",e[1]===e[4].code)},d(k){k&&v(t),r=!1,a()}}}function $e(n,e){let t,l,d,i;return l=new Oe({props:{content:e[4].body}}),{key:n,first:null,c(){t=p("div"),x(l.$$.fragment),d=y(),b(t,"class","tab-item"),Y(t,"active",e[1]===e[4].code),this.first=t},m(r,a){g(r,t,a),ee(l,t,null),u(t,d),i=!0},p(r,a){e=r;const m={};a&4&&(m.content=e[4].body),l.$set(m),(!i||a&6)&&Y(t,"active",e[1]===e[4].code)},i(r){i||(Q(l.$$.fragment,r),i=!0)},o(r){V(l.$$.fragment,r),i=!1},d(r){r&&v(t),te(l)}}}function Ne(n){let e,t,l,d,i,r,a,m=n[0].name+"",k,P,G,H,F,L,z,B,D,S,N,T=[],O=new Map,A,j,q=[],W=new Map,w,E=K(n[2]);const M=c=>c[4].code;for(let c=0;c<E.length;c+=1){let f=ve(n,E,c),s=M(f);O.set(s,T[c]=ke(s,f))}let _=K(n[2]);const X=c=>c[4].code;for(let c=0;c<_.length;c+=1){let f=ge(n,_,c),s=X(f);W.set(s,q[c]=$e(s,f))}return{c(){e=p("div"),t=p("strong"),t.textContent="POST",l=y(),d=p("div"),i=p("p"),r=U("/api/collections/"),a=p("strong"),k=U(m),P=U("/confirm-email-change"),G=y(),H=p("div"),H.textContent="Body Parameters",F=y(),L=p("table"),L.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>token</span></div></td> <td><span class="label">String</span></td> <td>The token from the change email request email.</td></tr> <tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>The account password to confirm the email change.</td></tr></tbody>',z=y(),B=p("div"),B.textContent="Responses",D=y(),S=p("div"),N=p("div");for(let c=0;c<T.length;c+=1)T[c].c();A=y(),j=p("div");for(let c=0;c<q.length;c+=1)q[c].c();b(t,"class","label label-primary"),b(d,"class","content"),b(e,"class","alert alert-success"),b(H,"class","section-title"),b(L,"class","table-compact table-border m-b-base"),b(B,"class","section-title"),b(N,"class","tabs-header compact combined left"),b(j,"class","tabs-content"),b(S,"class","tabs")},m(c,f){g(c,e,f),u(e,t),u(e,l),u(e,d),u(d,i),u(i,r),u(i,a),u(a,k),u(i,P),g(c,G,f),g(c,H,f),g(c,F,f),g(c,L,f),g(c,z,f),g(c,B,f),g(c,D,f),g(c,S,f),u(S,N);for(let s=0;s<T.length;s+=1)T[s]&&T[s].m(N,null);u(S,A),u(S,j);for(let s=0;s<q.length;s+=1)q[s]&&q[s].m(j,null);w=!0},p(c,[f]){(!w||f&1)&&m!==(m=c[0].name+"")&&J(k,m),f&6&&(E=K(c[2]),T=le(T,f,M,1,c,E,O,N,Re,ke,null,ve)),f&6&&(_=K(c[2]),ne(),q=le(q,f,X,1,c,_,W,j,Se,$e,null,ge),ae())},i(c){if(!w){for(let f=0;f<_.length;f+=1)Q(q[f]);w=!0}},o(c){for(let f=0;f<q.length;f+=1)V(q[f]);w=!1},d(c){c&&(v(e),v(G),v(H),v(F),v(L),v(z),v(B),v(D),v(S));for(let f=0;f<T.length;f+=1)T[f].d();for(let f=0;f<q.length;f+=1)q[f].d()}}}function We(n,e,t){let{collection:l}=e,d=204,i=[];const r=a=>t(1,d=a.code);return n.$$set=a=>{"collection"in a&&t(0,l=a.collection)},t(2,i=[{code:204,body:"null"},{code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to authenticate.",
                  "data": {
                    "token": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `}]),[l,d,i,r]}class He extends se{constructor(e){super(),oe(this,e,We,Ne,ie,{collection:0})}}function we(n,e,t){const l=n.slice();return l[4]=e[t],l}function Ce(n,e,t){const l=n.slice();return l[4]=e[t],l}function ye(n,e){let t,l=e[4].code+"",d,i,r,a;function m(){return e[3](e[4])}return{key:n,first:null,c(){t=p("button"),d=U(l),i=y(),b(t,"class","tab-item"),Y(t,"active",e[1]===e[4].code),this.first=t},m(k,P){g(k,t,P),u(t,d),u(t,i),r||(a=ce(t,"click",m),r=!0)},p(k,P){e=k,P&4&&l!==(l=e[4].code+"")&&J(d,l),P&6&&Y(t,"active",e[1]===e[4].code)},d(k){k&&v(t),r=!1,a()}}}function Ee(n,e){let t,l,d,i;return l=new Oe({props:{content:e[4].body}}),{key:n,first:null,c(){t=p("div"),x(l.$$.fragment),d=y(),b(t,"class","tab-item"),Y(t,"active",e[1]===e[4].code),this.first=t},m(r,a){g(r,t,a),ee(l,t,null),u(t,d),i=!0},p(r,a){e=r;const m={};a&4&&(m.content=e[4].body),l.$set(m),(!i||a&6)&&Y(t,"active",e[1]===e[4].code)},i(r){i||(Q(l.$$.fragment,r),i=!0)},o(r){V(l.$$.fragment,r),i=!1},d(r){r&&v(t),te(l)}}}function Le(n){let e,t,l,d,i,r,a,m=n[0].name+"",k,P,G,H,F,L,z,B,D,S,N,T,O,A=[],j=new Map,q,W,w=[],E=new Map,M,_=K(n[2]);const X=s=>s[4].code;for(let s=0;s<_.length;s+=1){let h=Ce(n,_,s),R=X(h);j.set(R,A[s]=ye(R,h))}let c=K(n[2]);const f=s=>s[4].code;for(let s=0;s<c.length;s+=1){let h=we(n,c,s),R=f(h);E.set(R,w[s]=Ee(R,h))}return{c(){e=p("div"),t=p("strong"),t.textContent="POST",l=y(),d=p("div"),i=p("p"),r=U("/api/collections/"),a=p("strong"),k=U(m),P=U("/request-email-change"),G=y(),H=p("p"),H.innerHTML="Requires <code>Authorization:TOKEN</code> header",F=y(),L=p("div"),L.textContent="Body Parameters",z=y(),B=p("table"),B.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>newEmail</span></div></td> <td><span class="label">String</span></td> <td>The new email address to send the change email request.</td></tr></tbody>',D=y(),S=p("div"),S.textContent="Responses",N=y(),T=p("div"),O=p("div");for(let s=0;s<A.length;s+=1)A[s].c();q=y(),W=p("div");for(let s=0;s<w.length;s+=1)w[s].c();b(t,"class","label label-primary"),b(d,"class","content"),b(H,"class","txt-hint txt-sm txt-right"),b(e,"class","alert alert-success"),b(L,"class","section-title"),b(B,"class","table-compact table-border m-b-base"),b(S,"class","section-title"),b(O,"class","tabs-header compact combined left"),b(W,"class","tabs-content"),b(T,"class","tabs")},m(s,h){g(s,e,h),u(e,t),u(e,l),u(e,d),u(d,i),u(i,r),u(i,a),u(a,k),u(i,P),u(e,G),u(e,H),g(s,F,h),g(s,L,h),g(s,z,h),g(s,B,h),g(s,D,h),g(s,S,h),g(s,N,h),g(s,T,h),u(T,O);for(let R=0;R<A.length;R+=1)A[R]&&A[R].m(O,null);u(T,q),u(T,W);for(let R=0;R<w.length;R+=1)w[R]&&w[R].m(W,null);M=!0},p(s,[h]){(!M||h&1)&&m!==(m=s[0].name+"")&&J(k,m),h&6&&(_=K(s[2]),A=le(A,h,X,1,s,_,j,O,Re,ye,null,Ce)),h&6&&(c=K(s[2]),ne(),w=le(w,h,f,1,s,c,E,W,Se,Ee,null,we),ae())},i(s){if(!M){for(let h=0;h<c.length;h+=1)Q(w[h]);M=!0}},o(s){for(let h=0;h<w.length;h+=1)V(w[h]);M=!1},d(s){s&&(v(e),v(F),v(L),v(z),v(B),v(D),v(S),v(N),v(T));for(let h=0;h<A.length;h+=1)A[h].d();for(let h=0;h<w.length;h+=1)w[h].d()}}}function Ue(n,e,t){let{collection:l}=e,d=204,i=[];const r=a=>t(1,d=a.code);return n.$$set=a=>{"collection"in a&&t(0,l=a.collection)},t(2,i=[{code:204,body:"null"},{code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to authenticate.",
                  "data": {
                    "newEmail": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `},{code:401,body:`
                {
                  "code": 401,
                  "message": "The request requires valid record authorization token to be set.",
                  "data": {}
                }
            `},{code:403,body:`
                {
                  "code": 403,
                  "message": "The authorized record model is not allowed to perform this action.",
                  "data": {}
                }
            `}]),[l,d,i,r]}class Ie extends se{constructor(e){super(),oe(this,e,Ue,Le,ie,{collection:0})}}function Te(n,e,t){const l=n.slice();return l[5]=e[t],l[7]=t,l}function qe(n,e,t){const l=n.slice();return l[5]=e[t],l[7]=t,l}function Pe(n){let e,t,l,d,i;function r(){return n[4](n[7])}return{c(){e=p("button"),t=p("div"),t.textContent=`${n[5].title}`,l=y(),b(t,"class","txt"),b(e,"class","tab-item"),Y(e,"active",n[1]==n[7])},m(a,m){g(a,e,m),u(e,t),u(e,l),d||(i=ce(e,"click",r),d=!0)},p(a,m){n=a,m&2&&Y(e,"active",n[1]==n[7])},d(a){a&&v(e),d=!1,i()}}}function Ae(n){let e,t,l,d;var i=n[5].component;function r(a,m){return{props:{collection:a[0]}}}return i&&(t=be(i,r(n))),{c(){e=p("div"),t&&x(t.$$.fragment),l=y(),b(e,"class","tab-item"),Y(e,"active",n[1]==n[7])},m(a,m){g(a,e,m),t&&ee(t,e,null),u(e,l),d=!0},p(a,m){if(i!==(i=a[5].component)){if(t){ne();const k=t;V(k.$$.fragment,1,0,()=>{te(k,1)}),ae()}i?(t=be(i,r(a)),x(t.$$.fragment),Q(t.$$.fragment,1),ee(t,e,l)):t=null}else if(i){const k={};m&1&&(k.collection=a[0]),t.$set(k)}(!d||m&2)&&Y(e,"active",a[1]==a[7])},i(a){d||(t&&Q(t.$$.fragment,a),d=!0)},o(a){t&&V(t.$$.fragment,a),d=!1},d(a){a&&v(e),t&&te(t)}}}function Ke(n){var c,f,s,h,R,re;let e,t,l=n[0].name+"",d,i,r,a,m,k,P,G=n[0].name+"",H,F,L,z,B,D,S,N,T,O,A,j,q,W;D=new Me({props:{js:`
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
    `}});let w=K(n[3]),E=[];for(let o=0;o<w.length;o+=1)E[o]=Pe(qe(n,w,o));let M=K(n[3]),_=[];for(let o=0;o<M.length;o+=1)_[o]=Ae(Te(n,M,o));const X=o=>V(_[o],1,1,()=>{_[o]=null});return{c(){e=p("h3"),t=U("Email change ("),d=U(l),i=U(")"),r=y(),a=p("div"),m=p("p"),k=U("Sends "),P=p("strong"),H=U(G),F=U(" email change request."),L=y(),z=p("p"),z.textContent=`On successful email change all previously issued auth tokens for the specific record will be
        automatically invalidated.`,B=y(),x(D.$$.fragment),S=y(),N=p("h6"),N.textContent="API details",T=y(),O=p("div"),A=p("div");for(let o=0;o<E.length;o+=1)E[o].c();j=y(),q=p("div");for(let o=0;o<_.length;o+=1)_[o].c();b(e,"class","m-b-sm"),b(a,"class","content txt-lg m-b-sm"),b(N,"class","m-b-xs"),b(A,"class","tabs-header compact"),b(q,"class","tabs-content"),b(O,"class","tabs")},m(o,C){g(o,e,C),u(e,t),u(e,d),u(e,i),g(o,r,C),g(o,a,C),u(a,m),u(m,k),u(m,P),u(P,H),u(m,F),u(a,L),u(a,z),g(o,B,C),ee(D,o,C),g(o,S,C),g(o,N,C),g(o,T,C),g(o,O,C),u(O,A);for(let I=0;I<E.length;I+=1)E[I]&&E[I].m(A,null);u(O,j),u(O,q);for(let I=0;I<_.length;I+=1)_[I]&&_[I].m(q,null);W=!0},p(o,[C]){var de,ue,fe,me,he,pe;(!W||C&1)&&l!==(l=o[0].name+"")&&J(d,l),(!W||C&1)&&G!==(G=o[0].name+"")&&J(H,G);const I={};if(C&5&&(I.js=`
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
    `),D.$set(I),C&10){w=K(o[3]);let $;for($=0;$<w.length;$+=1){const Z=qe(o,w,$);E[$]?E[$].p(Z,C):(E[$]=Pe(Z),E[$].c(),E[$].m(A,null))}for(;$<E.length;$+=1)E[$].d(1);E.length=w.length}if(C&11){M=K(o[3]);let $;for($=0;$<M.length;$+=1){const Z=Te(o,M,$);_[$]?(_[$].p(Z,C),Q(_[$],1)):(_[$]=Ae(Z),_[$].c(),Q(_[$],1),_[$].m(q,null))}for(ne(),$=M.length;$<_.length;$+=1)X($);ae()}},i(o){if(!W){Q(D.$$.fragment,o);for(let C=0;C<M.length;C+=1)Q(_[C]);W=!0}},o(o){V(D.$$.fragment,o),_=_.filter(Boolean);for(let C=0;C<_.length;C+=1)V(_[C]);W=!1},d(o){o&&(v(e),v(r),v(a),v(B),v(S),v(N),v(T),v(O)),te(D,o),_e(E,o),_e(_,o)}}}function Ye(n,e,t){let l,{collection:d}=e;const i=[{title:"Request email change",component:Ie},{title:"Confirm email change",component:He}];let r=0;const a=m=>t(1,r=m);return n.$$set=m=>{"collection"in m&&t(0,d=m.collection)},t(2,l=Be.getApiExampleUrl(De.baseURL)),[d,r,l,i,a]}class ze extends se{constructor(e){super(),oe(this,e,Ye,Ke,ie,{collection:0})}}export{ze as default};
