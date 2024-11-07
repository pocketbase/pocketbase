import{S as se,i as oe,s as ie,W as G,f as p,h as C,y as U,j as b,l as g,n as u,G as J,X as le,Y as Re,D as ne,Z as Se,E as ae,t as X,a as Z,u as v,k as K,o as ce,V as Oe,c as x,m as ee,d as te,U as Me,_ as _e,I as Be,p as De,$ as be}from"./index-D8ZreWy5.js";function ge(n,e,t){const l=n.slice();return l[4]=e[t],l}function ve(n,e,t){const l=n.slice();return l[4]=e[t],l}function ke(n,e){let t,l=e[4].code+"",d,i,r,a;function m(){return e[3](e[4])}return{key:n,first:null,c(){t=p("button"),d=U(l),i=C(),b(t,"class","tab-item"),K(t,"active",e[1]===e[4].code),this.first=t},m(k,A){g(k,t,A),u(t,d),u(t,i),r||(a=ce(t,"click",m),r=!0)},p(k,A){e=k,A&4&&l!==(l=e[4].code+"")&&J(d,l),A&6&&K(t,"active",e[1]===e[4].code)},d(k){k&&v(t),r=!1,a()}}}function $e(n,e){let t,l,d,i;return l=new Oe({props:{content:e[4].body}}),{key:n,first:null,c(){t=p("div"),x(l.$$.fragment),d=C(),b(t,"class","tab-item"),K(t,"active",e[1]===e[4].code),this.first=t},m(r,a){g(r,t,a),ee(l,t,null),u(t,d),i=!0},p(r,a){e=r;const m={};a&4&&(m.content=e[4].body),l.$set(m),(!i||a&6)&&K(t,"active",e[1]===e[4].code)},i(r){i||(X(l.$$.fragment,r),i=!0)},o(r){Z(l.$$.fragment,r),i=!1},d(r){r&&v(t),te(l)}}}function Ne(n){let e,t,l,d,i,r,a,m=n[0].name+"",k,A,Y,H,V,L,j,B,D,S,N,T=[],O=new Map,P,z,q=[],W=new Map,w,E=G(n[2]);const M=c=>c[4].code;for(let c=0;c<E.length;c+=1){let f=ve(n,E,c),s=M(f);O.set(s,T[c]=ke(s,f))}let _=G(n[2]);const F=c=>c[4].code;for(let c=0;c<_.length;c+=1){let f=ge(n,_,c),s=F(f);W.set(s,q[c]=$e(s,f))}return{c(){e=p("div"),t=p("strong"),t.textContent="POST",l=C(),d=p("div"),i=p("p"),r=U("/api/collections/"),a=p("strong"),k=U(m),A=U("/confirm-email-change"),Y=C(),H=p("div"),H.textContent="Body Parameters",V=C(),L=p("table"),L.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>token</span></div></td> <td><span class="label">String</span></td> <td>The token from the change email request email.</td></tr> <tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>The account password to confirm the email change.</td></tr></tbody>',j=C(),B=p("div"),B.textContent="Responses",D=C(),S=p("div"),N=p("div");for(let c=0;c<T.length;c+=1)T[c].c();P=C(),z=p("div");for(let c=0;c<q.length;c+=1)q[c].c();b(t,"class","label label-primary"),b(d,"class","content"),b(e,"class","alert alert-success"),b(H,"class","section-title"),b(L,"class","table-compact table-border m-b-base"),b(B,"class","section-title"),b(N,"class","tabs-header compact combined left"),b(z,"class","tabs-content"),b(S,"class","tabs")},m(c,f){g(c,e,f),u(e,t),u(e,l),u(e,d),u(d,i),u(i,r),u(i,a),u(a,k),u(i,A),g(c,Y,f),g(c,H,f),g(c,V,f),g(c,L,f),g(c,j,f),g(c,B,f),g(c,D,f),g(c,S,f),u(S,N);for(let s=0;s<T.length;s+=1)T[s]&&T[s].m(N,null);u(S,P),u(S,z);for(let s=0;s<q.length;s+=1)q[s]&&q[s].m(z,null);w=!0},p(c,[f]){(!w||f&1)&&m!==(m=c[0].name+"")&&J(k,m),f&6&&(E=G(c[2]),T=le(T,f,M,1,c,E,O,N,Re,ke,null,ve)),f&6&&(_=G(c[2]),ne(),q=le(q,f,F,1,c,_,W,z,Se,$e,null,ge),ae())},i(c){if(!w){for(let f=0;f<_.length;f+=1)X(q[f]);w=!0}},o(c){for(let f=0;f<q.length;f+=1)Z(q[f]);w=!1},d(c){c&&(v(e),v(Y),v(H),v(V),v(L),v(j),v(B),v(D),v(S));for(let f=0;f<T.length;f+=1)T[f].d();for(let f=0;f<q.length;f+=1)q[f].d()}}}function We(n,e,t){let{collection:l}=e,d=204,i=[];const r=a=>t(1,d=a.code);return n.$$set=a=>{"collection"in a&&t(0,l=a.collection)},t(2,i=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[l,d,i,r]}class He extends se{constructor(e){super(),oe(this,e,We,Ne,ie,{collection:0})}}function we(n,e,t){const l=n.slice();return l[4]=e[t],l}function ye(n,e,t){const l=n.slice();return l[4]=e[t],l}function Ce(n,e){let t,l=e[4].code+"",d,i,r,a;function m(){return e[3](e[4])}return{key:n,first:null,c(){t=p("button"),d=U(l),i=C(),b(t,"class","tab-item"),K(t,"active",e[1]===e[4].code),this.first=t},m(k,A){g(k,t,A),u(t,d),u(t,i),r||(a=ce(t,"click",m),r=!0)},p(k,A){e=k,A&4&&l!==(l=e[4].code+"")&&J(d,l),A&6&&K(t,"active",e[1]===e[4].code)},d(k){k&&v(t),r=!1,a()}}}function Ee(n,e){let t,l,d,i;return l=new Oe({props:{content:e[4].body}}),{key:n,first:null,c(){t=p("div"),x(l.$$.fragment),d=C(),b(t,"class","tab-item"),K(t,"active",e[1]===e[4].code),this.first=t},m(r,a){g(r,t,a),ee(l,t,null),u(t,d),i=!0},p(r,a){e=r;const m={};a&4&&(m.content=e[4].body),l.$set(m),(!i||a&6)&&K(t,"active",e[1]===e[4].code)},i(r){i||(X(l.$$.fragment,r),i=!0)},o(r){Z(l.$$.fragment,r),i=!1},d(r){r&&v(t),te(l)}}}function Le(n){let e,t,l,d,i,r,a,m=n[0].name+"",k,A,Y,H,V,L,j,B,D,S,N,T,O,P=[],z=new Map,q,W,w=[],E=new Map,M,_=G(n[2]);const F=s=>s[4].code;for(let s=0;s<_.length;s+=1){let h=ye(n,_,s),R=F(h);z.set(R,P[s]=Ce(R,h))}let c=G(n[2]);const f=s=>s[4].code;for(let s=0;s<c.length;s+=1){let h=we(n,c,s),R=f(h);E.set(R,w[s]=Ee(R,h))}return{c(){e=p("div"),t=p("strong"),t.textContent="POST",l=C(),d=p("div"),i=p("p"),r=U("/api/collections/"),a=p("strong"),k=U(m),A=U("/request-email-change"),Y=C(),H=p("p"),H.innerHTML="Requires <code>Authorization:TOKEN</code>",V=C(),L=p("div"),L.textContent="Body Parameters",j=C(),B=p("table"),B.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>newEmail</span></div></td> <td><span class="label">String</span></td> <td>The new email address to send the change email request.</td></tr></tbody>',D=C(),S=p("div"),S.textContent="Responses",N=C(),T=p("div"),O=p("div");for(let s=0;s<P.length;s+=1)P[s].c();q=C(),W=p("div");for(let s=0;s<w.length;s+=1)w[s].c();b(t,"class","label label-primary"),b(d,"class","content"),b(H,"class","txt-hint txt-sm txt-right"),b(e,"class","alert alert-success"),b(L,"class","section-title"),b(B,"class","table-compact table-border m-b-base"),b(S,"class","section-title"),b(O,"class","tabs-header compact combined left"),b(W,"class","tabs-content"),b(T,"class","tabs")},m(s,h){g(s,e,h),u(e,t),u(e,l),u(e,d),u(d,i),u(i,r),u(i,a),u(a,k),u(i,A),u(e,Y),u(e,H),g(s,V,h),g(s,L,h),g(s,j,h),g(s,B,h),g(s,D,h),g(s,S,h),g(s,N,h),g(s,T,h),u(T,O);for(let R=0;R<P.length;R+=1)P[R]&&P[R].m(O,null);u(T,q),u(T,W);for(let R=0;R<w.length;R+=1)w[R]&&w[R].m(W,null);M=!0},p(s,[h]){(!M||h&1)&&m!==(m=s[0].name+"")&&J(k,m),h&6&&(_=G(s[2]),P=le(P,h,F,1,s,_,z,O,Re,Ce,null,ye)),h&6&&(c=G(s[2]),ne(),w=le(w,h,f,1,s,c,E,W,Se,Ee,null,we),ae())},i(s){if(!M){for(let h=0;h<c.length;h+=1)X(w[h]);M=!0}},o(s){for(let h=0;h<w.length;h+=1)Z(w[h]);M=!1},d(s){s&&(v(e),v(V),v(L),v(j),v(B),v(D),v(S),v(N),v(T));for(let h=0;h<P.length;h+=1)P[h].d();for(let h=0;h<w.length;h+=1)w[h].d()}}}function Ue(n,e,t){let{collection:l}=e,d=204,i=[];const r=a=>t(1,d=a.code);return n.$$set=a=>{"collection"in a&&t(0,l=a.collection)},t(2,i=[{code:204,body:"null"},{code:400,body:`
                {
                  "code": 400,
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
            `}]),[l,d,i,r]}class Ie extends se{constructor(e){super(),oe(this,e,Ue,Le,ie,{collection:0})}}function Te(n,e,t){const l=n.slice();return l[5]=e[t],l[7]=t,l}function qe(n,e,t){const l=n.slice();return l[5]=e[t],l[7]=t,l}function Ae(n){let e,t,l,d,i;function r(){return n[4](n[7])}return{c(){e=p("button"),t=p("div"),t.textContent=`${n[5].title}`,l=C(),b(t,"class","txt"),b(e,"class","tab-item"),K(e,"active",n[1]==n[7])},m(a,m){g(a,e,m),u(e,t),u(e,l),d||(i=ce(e,"click",r),d=!0)},p(a,m){n=a,m&2&&K(e,"active",n[1]==n[7])},d(a){a&&v(e),d=!1,i()}}}function Pe(n){let e,t,l,d;var i=n[5].component;function r(a,m){return{props:{collection:a[0]}}}return i&&(t=be(i,r(n))),{c(){e=p("div"),t&&x(t.$$.fragment),l=C(),b(e,"class","tab-item"),K(e,"active",n[1]==n[7])},m(a,m){g(a,e,m),t&&ee(t,e,null),u(e,l),d=!0},p(a,m){if(i!==(i=a[5].component)){if(t){ne();const k=t;Z(k.$$.fragment,1,0,()=>{te(k,1)}),ae()}i?(t=be(i,r(a)),x(t.$$.fragment),X(t.$$.fragment,1),ee(t,e,l)):t=null}else if(i){const k={};m&1&&(k.collection=a[0]),t.$set(k)}(!d||m&2)&&K(e,"active",a[1]==a[7])},i(a){d||(t&&X(t.$$.fragment,a),d=!0)},o(a){t&&Z(t.$$.fragment,a),d=!1},d(a){a&&v(e),t&&te(t)}}}function Ge(n){var c,f,s,h,R,re;let e,t,l=n[0].name+"",d,i,r,a,m,k,A,Y=n[0].name+"",H,V,L,j,B,D,S,N,T,O,P,z,q,W;D=new Me({props:{js:`
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
    `}});let w=G(n[3]),E=[];for(let o=0;o<w.length;o+=1)E[o]=Ae(qe(n,w,o));let M=G(n[3]),_=[];for(let o=0;o<M.length;o+=1)_[o]=Pe(Te(n,M,o));const F=o=>Z(_[o],1,1,()=>{_[o]=null});return{c(){e=p("h3"),t=U("Email change ("),d=U(l),i=U(")"),r=C(),a=p("div"),m=p("p"),k=U("Sends "),A=p("strong"),H=U(Y),V=U(" email change request."),L=C(),j=p("p"),j.textContent=`On successful email change all previously issued auth tokens for the specific record will be
        automatically invalidated.`,B=C(),x(D.$$.fragment),S=C(),N=p("h6"),N.textContent="API details",T=C(),O=p("div"),P=p("div");for(let o=0;o<E.length;o+=1)E[o].c();z=C(),q=p("div");for(let o=0;o<_.length;o+=1)_[o].c();b(e,"class","m-b-sm"),b(a,"class","content txt-lg m-b-sm"),b(N,"class","m-b-xs"),b(P,"class","tabs-header compact"),b(q,"class","tabs-content"),b(O,"class","tabs")},m(o,y){g(o,e,y),u(e,t),u(e,d),u(e,i),g(o,r,y),g(o,a,y),u(a,m),u(m,k),u(m,A),u(A,H),u(m,V),u(a,L),u(a,j),g(o,B,y),ee(D,o,y),g(o,S,y),g(o,N,y),g(o,T,y),g(o,O,y),u(O,P);for(let I=0;I<E.length;I+=1)E[I]&&E[I].m(P,null);u(O,z),u(O,q);for(let I=0;I<_.length;I+=1)_[I]&&_[I].m(q,null);W=!0},p(o,[y]){var de,ue,fe,me,he,pe;(!W||y&1)&&l!==(l=o[0].name+"")&&J(d,l),(!W||y&1)&&Y!==(Y=o[0].name+"")&&J(H,Y);const I={};if(y&5&&(I.js=`
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
    `),y&5&&(I.dart=`
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
    `),D.$set(I),y&10){w=G(o[3]);let $;for($=0;$<w.length;$+=1){const Q=qe(o,w,$);E[$]?E[$].p(Q,y):(E[$]=Ae(Q),E[$].c(),E[$].m(P,null))}for(;$<E.length;$+=1)E[$].d(1);E.length=w.length}if(y&11){M=G(o[3]);let $;for($=0;$<M.length;$+=1){const Q=Te(o,M,$);_[$]?(_[$].p(Q,y),X(_[$],1)):(_[$]=Pe(Q),_[$].c(),X(_[$],1),_[$].m(q,null))}for(ne(),$=M.length;$<_.length;$+=1)F($);ae()}},i(o){if(!W){X(D.$$.fragment,o);for(let y=0;y<M.length;y+=1)X(_[y]);W=!0}},o(o){Z(D.$$.fragment,o),_=_.filter(Boolean);for(let y=0;y<_.length;y+=1)Z(_[y]);W=!1},d(o){o&&(v(e),v(r),v(a),v(B),v(S),v(N),v(T),v(O)),te(D,o),_e(E,o),_e(_,o)}}}function Ke(n,e,t){let l,{collection:d}=e;const i=[{title:"Request email change",component:Ie},{title:"Confirm email change",component:He}];let r=0;const a=m=>t(1,r=m);return n.$$set=m=>{"collection"in m&&t(0,d=m.collection)},t(2,l=Be.getApiExampleUrl(De.baseURL)),[d,r,l,i,a]}class je extends se{constructor(e){super(),oe(this,e,Ke,Ge,ie,{collection:0})}}export{je as default};
