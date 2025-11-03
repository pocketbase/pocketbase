import{S as be,i as _e,s as ve,ac as ge,O as J,h as b,d as x,t as Q,a as N,I as ce,ad as de,ae as Qe,C as ue,af as je,D as he,l as _,n as s,m as ee,u as d,v as T,A as R,c as te,w as g,J as ke,k as V,o as $e,ab as Ke,Q as De,p as ze,a4 as Me}from"./index-B4ZsHsKR.js";import{F as Ge}from"./FieldsQueryParam-K1y4zYh0.js";function Be(o,t,e){const l=o.slice();return l[4]=t[e],l}function Ie(o,t,e){const l=o.slice();return l[4]=t[e],l}function We(o,t){let e,l=t[4].code+"",h,i,c,n;function m(){return t[3](t[4])}return{key:o,first:null,c(){e=d("button"),h=R(l),i=T(),g(e,"class","tab-item"),V(e,"active",t[1]===t[4].code),this.first=e},m(v,C){_(v,e,C),s(e,h),s(e,i),c||(n=$e(e,"click",m),c=!0)},p(v,C){t=v,C&4&&l!==(l=t[4].code+"")&&ce(h,l),C&6&&V(e,"active",t[1]===t[4].code)},d(v){v&&b(e),c=!1,n()}}}function Fe(o,t){let e,l,h,i;return l=new ge({props:{content:t[4].body}}),{key:o,first:null,c(){e=d("div"),te(l.$$.fragment),h=T(),g(e,"class","tab-item"),V(e,"active",t[1]===t[4].code),this.first=e},m(c,n){_(c,e,n),ee(l,e,null),s(e,h),i=!0},p(c,n){t=c;const m={};n&4&&(m.content=t[4].body),l.$set(m),(!i||n&6)&&V(e,"active",t[1]===t[4].code)},i(c){i||(N(l.$$.fragment,c),i=!0)},o(c){Q(l.$$.fragment,c),i=!1},d(c){c&&b(e),x(l)}}}function Xe(o){let t,e,l,h,i,c,n,m=o[0].name+"",v,C,F,B,I,D,j,M,U,y,O,q,k,L,Y,A,z,E,a,$,P,X,u,p,S,w,G,we,Te,Pe,pe,Oe,ye,le,fe,ae,me,Z,oe,K=[],Se=new Map,qe,ne,H=[],Ce=new Map,se;P=new ge({props:{content:"?expand=relField1,relField2.subRelField"}}),le=new Ge({props:{prefix:"record."}});let re=J(o[2]);const Ae=r=>r[4].code;for(let r=0;r<re.length;r+=1){let f=Ie(o,re,r),W=Ae(f);Se.set(W,K[r]=We(W,f))}let ie=J(o[2]);const Re=r=>r[4].code;for(let r=0;r<ie.length;r+=1){let f=Be(o,ie,r),W=Re(f);Ce.set(W,H[r]=Fe(W,f))}return{c(){t=d("div"),e=d("strong"),e.textContent="POST",l=T(),h=d("div"),i=d("p"),c=R("/api/collections/"),n=d("strong"),v=R(m),C=R("/auth-with-otp"),F=T(),B=d("div"),B.textContent="Body Parameters",I=T(),D=d("table"),D.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>otpId</span></div></td> <td><span class="label">String</span></td> <td>The id of the OTP request.</td></tr> <tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>The one-time password.</td></tr></tbody>',j=T(),M=d("div"),M.textContent="Query parameters",U=T(),y=d("table"),O=d("thead"),O.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',q=T(),k=d("tbody"),L=d("tr"),Y=d("td"),Y.textContent="expand",A=T(),z=d("td"),z.innerHTML='<span class="label">String</span>',E=T(),a=d("td"),$=R(`Auto expand record relations. Ex.:
                `),te(P.$$.fragment),X=R(`
                Supports up to 6-levels depth nested relations expansion. `),u=d("br"),p=R(`
                The expanded relations will be appended to the record under the
                `),S=d("code"),S.textContent="expand",w=R(" property (eg. "),G=d("code"),G.textContent='"expand": {"relField1": {...}, ...}',we=R(`).
                `),Te=d("br"),Pe=R(`
                Only the relations to which the request user has permissions to `),pe=d("strong"),pe.textContent="view",Oe=R(" will be expanded."),ye=T(),te(le.$$.fragment),fe=T(),ae=d("div"),ae.textContent="Responses",me=T(),Z=d("div"),oe=d("div");for(let r=0;r<K.length;r+=1)K[r].c();qe=T(),ne=d("div");for(let r=0;r<H.length;r+=1)H[r].c();g(e,"class","label label-primary"),g(h,"class","content"),g(t,"class","alert alert-success"),g(B,"class","section-title"),g(D,"class","table-compact table-border m-b-base"),g(M,"class","section-title"),g(y,"class","table-compact table-border m-b-base"),g(ae,"class","section-title"),g(oe,"class","tabs-header compact combined left"),g(ne,"class","tabs-content"),g(Z,"class","tabs")},m(r,f){_(r,t,f),s(t,e),s(t,l),s(t,h),s(h,i),s(i,c),s(i,n),s(n,v),s(i,C),_(r,F,f),_(r,B,f),_(r,I,f),_(r,D,f),_(r,j,f),_(r,M,f),_(r,U,f),_(r,y,f),s(y,O),s(y,q),s(y,k),s(k,L),s(L,Y),s(L,A),s(L,z),s(L,E),s(L,a),s(a,$),ee(P,a,null),s(a,X),s(a,u),s(a,p),s(a,S),s(a,w),s(a,G),s(a,we),s(a,Te),s(a,Pe),s(a,pe),s(a,Oe),s(k,ye),ee(le,k,null),_(r,fe,f),_(r,ae,f),_(r,me,f),_(r,Z,f),s(Z,oe);for(let W=0;W<K.length;W+=1)K[W]&&K[W].m(oe,null);s(Z,qe),s(Z,ne);for(let W=0;W<H.length;W+=1)H[W]&&H[W].m(ne,null);se=!0},p(r,[f]){(!se||f&1)&&m!==(m=r[0].name+"")&&ce(v,m),f&6&&(re=J(r[2]),K=de(K,f,Ae,1,r,re,Se,oe,Qe,We,null,Ie)),f&6&&(ie=J(r[2]),ue(),H=de(H,f,Re,1,r,ie,Ce,ne,je,Fe,null,Be),he())},i(r){if(!se){N(P.$$.fragment,r),N(le.$$.fragment,r);for(let f=0;f<ie.length;f+=1)N(H[f]);se=!0}},o(r){Q(P.$$.fragment,r),Q(le.$$.fragment,r);for(let f=0;f<H.length;f+=1)Q(H[f]);se=!1},d(r){r&&(b(t),b(F),b(B),b(I),b(D),b(j),b(M),b(U),b(y),b(fe),b(ae),b(me),b(Z)),x(P),x(le);for(let f=0;f<K.length;f+=1)K[f].d();for(let f=0;f<H.length;f+=1)H[f].d()}}}function Ze(o,t,e){let{collection:l}=t,h=200,i=[];const c=n=>e(1,h=n.code);return o.$$set=n=>{"collection"in n&&e(0,l=n.collection)},o.$$.update=()=>{o.$$.dirty&1&&e(2,i=[{code:200,body:JSON.stringify({token:"JWT_TOKEN",record:ke.dummyCollectionRecord(l)},null,2)},{code:400,body:`
                {
                  "status": 400,
                  "message": "Failed to authenticate.",
                  "data": {
                    "otpId": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `}])},[l,h,i,c]}class xe extends be{constructor(t){super(),_e(this,t,Ze,Xe,ve,{collection:0})}}function Ue(o,t,e){const l=o.slice();return l[4]=t[e],l}function He(o,t,e){const l=o.slice();return l[4]=t[e],l}function Le(o,t){let e,l=t[4].code+"",h,i,c,n;function m(){return t[3](t[4])}return{key:o,first:null,c(){e=d("button"),h=R(l),i=T(),g(e,"class","tab-item"),V(e,"active",t[1]===t[4].code),this.first=e},m(v,C){_(v,e,C),s(e,h),s(e,i),c||(n=$e(e,"click",m),c=!0)},p(v,C){t=v,C&4&&l!==(l=t[4].code+"")&&ce(h,l),C&6&&V(e,"active",t[1]===t[4].code)},d(v){v&&b(e),c=!1,n()}}}function Ye(o,t){let e,l,h,i;return l=new ge({props:{content:t[4].body}}),{key:o,first:null,c(){e=d("div"),te(l.$$.fragment),h=T(),g(e,"class","tab-item"),V(e,"active",t[1]===t[4].code),this.first=e},m(c,n){_(c,e,n),ee(l,e,null),s(e,h),i=!0},p(c,n){t=c;const m={};n&4&&(m.content=t[4].body),l.$set(m),(!i||n&6)&&V(e,"active",t[1]===t[4].code)},i(c){i||(N(l.$$.fragment,c),i=!0)},o(c){Q(l.$$.fragment,c),i=!1},d(c){c&&b(e),x(l)}}}function et(o){let t,e,l,h,i,c,n,m=o[0].name+"",v,C,F,B,I,D,j,M,U,y,O,q=[],k=new Map,L,Y,A=[],z=new Map,E,a=J(o[2]);const $=u=>u[4].code;for(let u=0;u<a.length;u+=1){let p=He(o,a,u),S=$(p);k.set(S,q[u]=Le(S,p))}let P=J(o[2]);const X=u=>u[4].code;for(let u=0;u<P.length;u+=1){let p=Ue(o,P,u),S=X(p);z.set(S,A[u]=Ye(S,p))}return{c(){t=d("div"),e=d("strong"),e.textContent="POST",l=T(),h=d("div"),i=d("p"),c=R("/api/collections/"),n=d("strong"),v=R(m),C=R("/request-otp"),F=T(),B=d("div"),B.textContent="Body Parameters",I=T(),D=d("table"),D.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>email</span></div></td> <td><span class="label">String</span></td> <td>The auth record email address to send the OTP request (if exists).</td></tr></tbody>',j=T(),M=d("div"),M.textContent="Responses",U=T(),y=d("div"),O=d("div");for(let u=0;u<q.length;u+=1)q[u].c();L=T(),Y=d("div");for(let u=0;u<A.length;u+=1)A[u].c();g(e,"class","label label-primary"),g(h,"class","content"),g(t,"class","alert alert-success"),g(B,"class","section-title"),g(D,"class","table-compact table-border m-b-base"),g(M,"class","section-title"),g(O,"class","tabs-header compact combined left"),g(Y,"class","tabs-content"),g(y,"class","tabs")},m(u,p){_(u,t,p),s(t,e),s(t,l),s(t,h),s(h,i),s(i,c),s(i,n),s(n,v),s(i,C),_(u,F,p),_(u,B,p),_(u,I,p),_(u,D,p),_(u,j,p),_(u,M,p),_(u,U,p),_(u,y,p),s(y,O);for(let S=0;S<q.length;S+=1)q[S]&&q[S].m(O,null);s(y,L),s(y,Y);for(let S=0;S<A.length;S+=1)A[S]&&A[S].m(Y,null);E=!0},p(u,[p]){(!E||p&1)&&m!==(m=u[0].name+"")&&ce(v,m),p&6&&(a=J(u[2]),q=de(q,p,$,1,u,a,k,O,Qe,Le,null,He)),p&6&&(P=J(u[2]),ue(),A=de(A,p,X,1,u,P,z,Y,je,Ye,null,Ue),he())},i(u){if(!E){for(let p=0;p<P.length;p+=1)N(A[p]);E=!0}},o(u){for(let p=0;p<A.length;p+=1)Q(A[p]);E=!1},d(u){u&&(b(t),b(F),b(B),b(I),b(D),b(j),b(M),b(U),b(y));for(let p=0;p<q.length;p+=1)q[p].d();for(let p=0;p<A.length;p+=1)A[p].d()}}}function tt(o,t,e){let{collection:l}=t,h=200,i=[];const c=n=>e(1,h=n.code);return o.$$set=n=>{"collection"in n&&e(0,l=n.collection)},e(2,i=[{code:200,body:JSON.stringify({otpId:ke.randomString(15)},null,2)},{code:400,body:`
                {
                  "status": 400,
                  "message": "An error occurred while validating the submitted data.",
                  "data": {
                    "email": {
                      "code": "validation_is_email",
                      "message": "Must be a valid email address."
                    }
                  }
                }
            `},{code:429,body:`
                {
                  "status": 429,
                  "message": "You've send too many OTP requests, please try again later.",
                  "data": {}
                }
            `}]),[l,h,i,c]}class lt extends be{constructor(t){super(),_e(this,t,tt,et,ve,{collection:0})}}function Je(o,t,e){const l=o.slice();return l[5]=t[e],l[7]=e,l}function Ne(o,t,e){const l=o.slice();return l[5]=t[e],l[7]=e,l}function Ve(o){let t,e,l,h,i;function c(){return o[4](o[7])}return{c(){t=d("button"),e=d("div"),e.textContent=`${o[5].title}`,l=T(),g(e,"class","txt"),g(t,"class","tab-item"),V(t,"active",o[1]==o[7])},m(n,m){_(n,t,m),s(t,e),s(t,l),h||(i=$e(t,"click",c),h=!0)},p(n,m){o=n,m&2&&V(t,"active",o[1]==o[7])},d(n){n&&b(t),h=!1,i()}}}function Ee(o){let t,e,l,h;var i=o[5].component;function c(n,m){return{props:{collection:n[0]}}}return i&&(e=Me(i,c(o))),{c(){t=d("div"),e&&te(e.$$.fragment),l=T(),g(t,"class","tab-item"),V(t,"active",o[1]==o[7])},m(n,m){_(n,t,m),e&&ee(e,t,null),s(t,l),h=!0},p(n,m){if(i!==(i=n[5].component)){if(e){ue();const v=e;Q(v.$$.fragment,1,0,()=>{x(v,1)}),he()}i?(e=Me(i,c(n)),te(e.$$.fragment),N(e.$$.fragment,1),ee(e,t,l)):e=null}else if(i){const v={};m&1&&(v.collection=n[0]),e.$set(v)}(!h||m&2)&&V(t,"active",n[1]==n[7])},i(n){h||(e&&N(e.$$.fragment,n),h=!0)},o(n){e&&Q(e.$$.fragment,n),h=!1},d(n){n&&b(t),e&&x(e)}}}function at(o){var Y,A,z,E;let t,e,l=o[0].name+"",h,i,c,n,m,v,C,F,B,I,D,j,M,U;v=new Ke({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${o[2]}');

        ...

        // send OTP email to the provided auth record
        const req = await pb.collection('${(Y=o[0])==null?void 0:Y.name}').requestOTP('test@example.com');

        // ... show a screen/popup to enter the password from the email ...

        // authenticate with the requested OTP id and the email password
        const authData = await pb.collection('${(A=o[0])==null?void 0:A.name}').authWithOTP(
            req.otpId,
            "YOUR_OTP",
        );

        // after the above you can also access the auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.record.id);

        // "logout"
        pb.authStore.clear();
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${o[2]}');

        ...

        // send OTP email to the provided auth record
        final req = await pb.collection('${(z=o[0])==null?void 0:z.name}').requestOTP('test@example.com');

        // ... show a screen/popup to enter the password from the email ...

        // authenticate with the requested OTP id and the email password
        final authData = await pb.collection('${(E=o[0])==null?void 0:E.name}').authWithOTP(
            req.otpId,
            "YOUR_OTP",
        );

        // after the above you can also access the auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.record.id);

        // "logout"
        pb.authStore.clear();
    `}});let y=J(o[3]),O=[];for(let a=0;a<y.length;a+=1)O[a]=Ve(Ne(o,y,a));let q=J(o[3]),k=[];for(let a=0;a<q.length;a+=1)k[a]=Ee(Je(o,q,a));const L=a=>Q(k[a],1,1,()=>{k[a]=null});return{c(){t=d("h3"),e=R("Auth with OTP ("),h=R(l),i=R(")"),c=T(),n=d("div"),n.innerHTML=`<p>Authenticate with an one-time password (OTP).</p> <p>Note that when requesting an OTP we return an <code>otpId</code> even if a user with the provided email
        doesn&#39;t exist as a very basic enumeration protection.</p>`,m=T(),te(v.$$.fragment),C=T(),F=d("h6"),F.textContent="API details",B=T(),I=d("div"),D=d("div");for(let a=0;a<O.length;a+=1)O[a].c();j=T(),M=d("div");for(let a=0;a<k.length;a+=1)k[a].c();g(t,"class","m-b-sm"),g(n,"class","content txt-lg m-b-sm"),g(F,"class","m-b-xs"),g(D,"class","tabs-header compact"),g(M,"class","tabs-content"),g(I,"class","tabs")},m(a,$){_(a,t,$),s(t,e),s(t,h),s(t,i),_(a,c,$),_(a,n,$),_(a,m,$),ee(v,a,$),_(a,C,$),_(a,F,$),_(a,B,$),_(a,I,$),s(I,D);for(let P=0;P<O.length;P+=1)O[P]&&O[P].m(D,null);s(I,j),s(I,M);for(let P=0;P<k.length;P+=1)k[P]&&k[P].m(M,null);U=!0},p(a,[$]){var X,u,p,S;(!U||$&1)&&l!==(l=a[0].name+"")&&ce(h,l);const P={};if($&5&&(P.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${a[2]}');

        ...

        // send OTP email to the provided auth record
        const req = await pb.collection('${(X=a[0])==null?void 0:X.name}').requestOTP('test@example.com');

        // ... show a screen/popup to enter the password from the email ...

        // authenticate with the requested OTP id and the email password
        const authData = await pb.collection('${(u=a[0])==null?void 0:u.name}').authWithOTP(
            req.otpId,
            "YOUR_OTP",
        );

        // after the above you can also access the auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.record.id);

        // "logout"
        pb.authStore.clear();
    `),$&5&&(P.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${a[2]}');

        ...

        // send OTP email to the provided auth record
        final req = await pb.collection('${(p=a[0])==null?void 0:p.name}').requestOTP('test@example.com');

        // ... show a screen/popup to enter the password from the email ...

        // authenticate with the requested OTP id and the email password
        final authData = await pb.collection('${(S=a[0])==null?void 0:S.name}').authWithOTP(
            req.otpId,
            "YOUR_OTP",
        );

        // after the above you can also access the auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.record.id);

        // "logout"
        pb.authStore.clear();
    `),v.$set(P),$&10){y=J(a[3]);let w;for(w=0;w<y.length;w+=1){const G=Ne(a,y,w);O[w]?O[w].p(G,$):(O[w]=Ve(G),O[w].c(),O[w].m(D,null))}for(;w<O.length;w+=1)O[w].d(1);O.length=y.length}if($&11){q=J(a[3]);let w;for(w=0;w<q.length;w+=1){const G=Je(a,q,w);k[w]?(k[w].p(G,$),N(k[w],1)):(k[w]=Ee(G),k[w].c(),N(k[w],1),k[w].m(M,null))}for(ue(),w=q.length;w<k.length;w+=1)L(w);he()}},i(a){if(!U){N(v.$$.fragment,a);for(let $=0;$<q.length;$+=1)N(k[$]);U=!0}},o(a){Q(v.$$.fragment,a),k=k.filter(Boolean);for(let $=0;$<k.length;$+=1)Q(k[$]);U=!1},d(a){a&&(b(t),b(c),b(n),b(m),b(C),b(F),b(B),b(I)),x(v,a),De(O,a),De(k,a)}}}function ot(o,t,e){let l,{collection:h}=t;const i=[{title:"OTP Request",component:lt},{title:"OTP Auth",component:xe}];let c=0;const n=m=>e(1,c=m);return o.$$set=m=>{"collection"in m&&e(0,h=m.collection)},e(2,l=ke.getApiExampleUrl(ze.baseURL)),[h,c,l,i,n]}class it extends be{constructor(t){super(),_e(this,t,ot,at,ve,{collection:0})}}export{it as default};
