import{S as be,i as _e,s as ve,W as ge,X as V,h as d,j as T,z as R,c as x,k as g,n as b,o as s,m as ee,H as ce,Y as de,Z as je,E as ue,_ as Qe,G as he,t as E,a as j,v as _,d as te,J as ke,l as J,q as $e,V as ze,$ as Me,p as Ge,a0 as Be}from"./index-SKn09NMF.js";import{F as Ke}from"./FieldsQueryParam-NXRpAlxi.js";function De(a,t,e){const l=a.slice();return l[4]=t[e],l}function Ie(a,t,e){const l=a.slice();return l[4]=t[e],l}function We(a,t){let e,l=t[4].code+"",h,i,c,n;function m(){return t[3](t[4])}return{key:a,first:null,c(){e=d("button"),h=R(l),i=T(),g(e,"class","tab-item"),J(e,"active",t[1]===t[4].code),this.first=e},m(v,C){b(v,e,C),s(e,h),s(e,i),c||(n=$e(e,"click",m),c=!0)},p(v,C){t=v,C&4&&l!==(l=t[4].code+"")&&ce(h,l),C&6&&J(e,"active",t[1]===t[4].code)},d(v){v&&_(e),c=!1,n()}}}function Fe(a,t){let e,l,h,i;return l=new ge({props:{content:t[4].body}}),{key:a,first:null,c(){e=d("div"),x(l.$$.fragment),h=T(),g(e,"class","tab-item"),J(e,"active",t[1]===t[4].code),this.first=e},m(c,n){b(c,e,n),ee(l,e,null),s(e,h),i=!0},p(c,n){t=c;const m={};n&4&&(m.content=t[4].body),l.$set(m),(!i||n&6)&&J(e,"active",t[1]===t[4].code)},i(c){i||(E(l.$$.fragment,c),i=!0)},o(c){j(l.$$.fragment,c),i=!1},d(c){c&&_(e),te(l)}}}function Xe(a){let t,e,l,h,i,c,n,m=a[0].name+"",v,C,F,D,I,M,Q,B,H,y,O,q,k,L,Y,A,G,N,o,$,P,X,u,p,S,w,K,we,Te,Pe,pe,Oe,ye,le,fe,oe,me,Z,ae,z=[],Se=new Map,qe,ne,U=[],Ce=new Map,se;P=new ge({props:{content:"?expand=relField1,relField2.subRelField"}}),le=new Ke({props:{prefix:"record."}});let re=V(a[2]);const Ae=r=>r[4].code;for(let r=0;r<re.length;r+=1){let f=Ie(a,re,r),W=Ae(f);Se.set(W,z[r]=We(W,f))}let ie=V(a[2]);const Re=r=>r[4].code;for(let r=0;r<ie.length;r+=1){let f=De(a,ie,r),W=Re(f);Ce.set(W,U[r]=Fe(W,f))}return{c(){t=d("div"),e=d("strong"),e.textContent="POST",l=T(),h=d("div"),i=d("p"),c=R("/api/collections/"),n=d("strong"),v=R(m),C=R("/auth-with-otp"),F=T(),D=d("div"),D.textContent="Body Parameters",I=T(),M=d("table"),M.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>otpId</span></div></td> <td><span class="label">String</span></td> <td>The id of the OTP request.</td></tr> <tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>The one-time password.</td></tr></tbody>',Q=T(),B=d("div"),B.textContent="Query parameters",H=T(),y=d("table"),O=d("thead"),O.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',q=T(),k=d("tbody"),L=d("tr"),Y=d("td"),Y.textContent="expand",A=T(),G=d("td"),G.innerHTML='<span class="label">String</span>',N=T(),o=d("td"),$=R(`Auto expand record relations. Ex.:
                `),x(P.$$.fragment),X=R(`
                Supports up to 6-levels depth nested relations expansion. `),u=d("br"),p=R(`
                The expanded relations will be appended to the record under the
                `),S=d("code"),S.textContent="expand",w=R(" property (eg. "),K=d("code"),K.textContent='"expand": {"relField1": {...}, ...}',we=R(`).
                `),Te=d("br"),Pe=R(`
                Only the relations to which the request user has permissions to `),pe=d("strong"),pe.textContent="view",Oe=R(" will be expanded."),ye=T(),x(le.$$.fragment),fe=T(),oe=d("div"),oe.textContent="Responses",me=T(),Z=d("div"),ae=d("div");for(let r=0;r<z.length;r+=1)z[r].c();qe=T(),ne=d("div");for(let r=0;r<U.length;r+=1)U[r].c();g(e,"class","label label-primary"),g(h,"class","content"),g(t,"class","alert alert-success"),g(D,"class","section-title"),g(M,"class","table-compact table-border m-b-base"),g(B,"class","section-title"),g(y,"class","table-compact table-border m-b-base"),g(oe,"class","section-title"),g(ae,"class","tabs-header compact combined left"),g(ne,"class","tabs-content"),g(Z,"class","tabs")},m(r,f){b(r,t,f),s(t,e),s(t,l),s(t,h),s(h,i),s(i,c),s(i,n),s(n,v),s(i,C),b(r,F,f),b(r,D,f),b(r,I,f),b(r,M,f),b(r,Q,f),b(r,B,f),b(r,H,f),b(r,y,f),s(y,O),s(y,q),s(y,k),s(k,L),s(L,Y),s(L,A),s(L,G),s(L,N),s(L,o),s(o,$),ee(P,o,null),s(o,X),s(o,u),s(o,p),s(o,S),s(o,w),s(o,K),s(o,we),s(o,Te),s(o,Pe),s(o,pe),s(o,Oe),s(k,ye),ee(le,k,null),b(r,fe,f),b(r,oe,f),b(r,me,f),b(r,Z,f),s(Z,ae);for(let W=0;W<z.length;W+=1)z[W]&&z[W].m(ae,null);s(Z,qe),s(Z,ne);for(let W=0;W<U.length;W+=1)U[W]&&U[W].m(ne,null);se=!0},p(r,[f]){(!se||f&1)&&m!==(m=r[0].name+"")&&ce(v,m),f&6&&(re=V(r[2]),z=de(z,f,Ae,1,r,re,Se,ae,je,We,null,Ie)),f&6&&(ie=V(r[2]),ue(),U=de(U,f,Re,1,r,ie,Ce,ne,Qe,Fe,null,De),he())},i(r){if(!se){E(P.$$.fragment,r),E(le.$$.fragment,r);for(let f=0;f<ie.length;f+=1)E(U[f]);se=!0}},o(r){j(P.$$.fragment,r),j(le.$$.fragment,r);for(let f=0;f<U.length;f+=1)j(U[f]);se=!1},d(r){r&&(_(t),_(F),_(D),_(I),_(M),_(Q),_(B),_(H),_(y),_(fe),_(oe),_(me),_(Z)),te(P),te(le);for(let f=0;f<z.length;f+=1)z[f].d();for(let f=0;f<U.length;f+=1)U[f].d()}}}function Ze(a,t,e){let{collection:l}=t,h=200,i=[];const c=n=>e(1,h=n.code);return a.$$set=n=>{"collection"in n&&e(0,l=n.collection)},a.$$.update=()=>{a.$$.dirty&1&&e(2,i=[{code:200,body:JSON.stringify({token:"JWT_TOKEN",record:ke.dummyCollectionRecord(l)},null,2)},{code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to authenticate.",
                  "data": {
                    "otpId": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `}])},[l,h,i,c]}class xe extends be{constructor(t){super(),_e(this,t,Ze,Xe,ve,{collection:0})}}function He(a,t,e){const l=a.slice();return l[4]=t[e],l}function Ue(a,t,e){const l=a.slice();return l[4]=t[e],l}function Le(a,t){let e,l=t[4].code+"",h,i,c,n;function m(){return t[3](t[4])}return{key:a,first:null,c(){e=d("button"),h=R(l),i=T(),g(e,"class","tab-item"),J(e,"active",t[1]===t[4].code),this.first=e},m(v,C){b(v,e,C),s(e,h),s(e,i),c||(n=$e(e,"click",m),c=!0)},p(v,C){t=v,C&4&&l!==(l=t[4].code+"")&&ce(h,l),C&6&&J(e,"active",t[1]===t[4].code)},d(v){v&&_(e),c=!1,n()}}}function Ye(a,t){let e,l,h,i;return l=new ge({props:{content:t[4].body}}),{key:a,first:null,c(){e=d("div"),x(l.$$.fragment),h=T(),g(e,"class","tab-item"),J(e,"active",t[1]===t[4].code),this.first=e},m(c,n){b(c,e,n),ee(l,e,null),s(e,h),i=!0},p(c,n){t=c;const m={};n&4&&(m.content=t[4].body),l.$set(m),(!i||n&6)&&J(e,"active",t[1]===t[4].code)},i(c){i||(E(l.$$.fragment,c),i=!0)},o(c){j(l.$$.fragment,c),i=!1},d(c){c&&_(e),te(l)}}}function et(a){let t,e,l,h,i,c,n,m=a[0].name+"",v,C,F,D,I,M,Q,B,H,y,O,q=[],k=new Map,L,Y,A=[],G=new Map,N,o=V(a[2]);const $=u=>u[4].code;for(let u=0;u<o.length;u+=1){let p=Ue(a,o,u),S=$(p);k.set(S,q[u]=Le(S,p))}let P=V(a[2]);const X=u=>u[4].code;for(let u=0;u<P.length;u+=1){let p=He(a,P,u),S=X(p);G.set(S,A[u]=Ye(S,p))}return{c(){t=d("div"),e=d("strong"),e.textContent="POST",l=T(),h=d("div"),i=d("p"),c=R("/api/collections/"),n=d("strong"),v=R(m),C=R("/request-otp"),F=T(),D=d("div"),D.textContent="Body Parameters",I=T(),M=d("table"),M.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>email</span></div></td> <td><span class="label">String</span></td> <td>The auth record email address to send the OTP request (if exists).</td></tr></tbody>',Q=T(),B=d("div"),B.textContent="Responses",H=T(),y=d("div"),O=d("div");for(let u=0;u<q.length;u+=1)q[u].c();L=T(),Y=d("div");for(let u=0;u<A.length;u+=1)A[u].c();g(e,"class","label label-primary"),g(h,"class","content"),g(t,"class","alert alert-success"),g(D,"class","section-title"),g(M,"class","table-compact table-border m-b-base"),g(B,"class","section-title"),g(O,"class","tabs-header compact combined left"),g(Y,"class","tabs-content"),g(y,"class","tabs")},m(u,p){b(u,t,p),s(t,e),s(t,l),s(t,h),s(h,i),s(i,c),s(i,n),s(n,v),s(i,C),b(u,F,p),b(u,D,p),b(u,I,p),b(u,M,p),b(u,Q,p),b(u,B,p),b(u,H,p),b(u,y,p),s(y,O);for(let S=0;S<q.length;S+=1)q[S]&&q[S].m(O,null);s(y,L),s(y,Y);for(let S=0;S<A.length;S+=1)A[S]&&A[S].m(Y,null);N=!0},p(u,[p]){(!N||p&1)&&m!==(m=u[0].name+"")&&ce(v,m),p&6&&(o=V(u[2]),q=de(q,p,$,1,u,o,k,O,je,Le,null,Ue)),p&6&&(P=V(u[2]),ue(),A=de(A,p,X,1,u,P,G,Y,Qe,Ye,null,He),he())},i(u){if(!N){for(let p=0;p<P.length;p+=1)E(A[p]);N=!0}},o(u){for(let p=0;p<A.length;p+=1)j(A[p]);N=!1},d(u){u&&(_(t),_(F),_(D),_(I),_(M),_(Q),_(B),_(H),_(y));for(let p=0;p<q.length;p+=1)q[p].d();for(let p=0;p<A.length;p+=1)A[p].d()}}}function tt(a,t,e){let{collection:l}=t,h=200,i=[];const c=n=>e(1,h=n.code);return a.$$set=n=>{"collection"in n&&e(0,l=n.collection)},e(2,i=[{code:200,body:JSON.stringify({otpId:ke.randomString(15)},null,2)},{code:400,body:`
                {
                  "code": 400,
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
                  "code": 429,
                  "message": "You've send too many OTP requests, please try again later.",
                  "data": {}
                }
            `}]),[l,h,i,c]}class lt extends be{constructor(t){super(),_e(this,t,tt,et,ve,{collection:0})}}function Ve(a,t,e){const l=a.slice();return l[5]=t[e],l[7]=e,l}function Ee(a,t,e){const l=a.slice();return l[5]=t[e],l[7]=e,l}function Je(a){let t,e,l,h,i;function c(){return a[4](a[7])}return{c(){t=d("button"),e=d("div"),e.textContent=`${a[5].title}`,l=T(),g(e,"class","txt"),g(t,"class","tab-item"),J(t,"active",a[1]==a[7])},m(n,m){b(n,t,m),s(t,e),s(t,l),h||(i=$e(t,"click",c),h=!0)},p(n,m){a=n,m&2&&J(t,"active",a[1]==a[7])},d(n){n&&_(t),h=!1,i()}}}function Ne(a){let t,e,l,h;var i=a[5].component;function c(n,m){return{props:{collection:n[0]}}}return i&&(e=Be(i,c(a))),{c(){t=d("div"),e&&x(e.$$.fragment),l=T(),g(t,"class","tab-item"),J(t,"active",a[1]==a[7])},m(n,m){b(n,t,m),e&&ee(e,t,null),s(t,l),h=!0},p(n,m){if(i!==(i=n[5].component)){if(e){ue();const v=e;j(v.$$.fragment,1,0,()=>{te(v,1)}),he()}i?(e=Be(i,c(n)),x(e.$$.fragment),E(e.$$.fragment,1),ee(e,t,l)):e=null}else if(i){const v={};m&1&&(v.collection=n[0]),e.$set(v)}(!h||m&2)&&J(t,"active",n[1]==n[7])},i(n){h||(e&&E(e.$$.fragment,n),h=!0)},o(n){e&&j(e.$$.fragment,n),h=!1},d(n){n&&_(t),e&&te(e)}}}function ot(a){var Y,A,G,N;let t,e,l=a[0].name+"",h,i,c,n,m,v,C,F,D,I,M,Q,B,H;v=new ze({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${a[2]}');

        ...

        // send OTP email to the provided auth record
        const req = await pb.collection('${(Y=a[0])==null?void 0:Y.name}').requestOTP('test@example.com');

        // ... show a screen/popup to enter the password from the email ...

        // authenticate with the requested OTP id and the email password
        const authData = await pb.collection('${(A=a[0])==null?void 0:A.name}').authWithOTP(
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

        final pb = PocketBase('${a[2]}');

        ...

        // send OTP email to the provided auth record
        final req = await pb.collection('${(G=a[0])==null?void 0:G.name}').requestOTP('test@example.com');

        // ... show a screen/popup to enter the password from the email ...

        // authenticate with the requested OTP id and the email password
        final authData = await pb.collection('${(N=a[0])==null?void 0:N.name}').authWithOTP(
            req.otpId,
            "YOUR_OTP",
        );

        // after the above you can also access the auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.record.id);

        // "logout"
        pb.authStore.clear();
    `}});let y=V(a[3]),O=[];for(let o=0;o<y.length;o+=1)O[o]=Je(Ee(a,y,o));let q=V(a[3]),k=[];for(let o=0;o<q.length;o+=1)k[o]=Ne(Ve(a,q,o));const L=o=>j(k[o],1,1,()=>{k[o]=null});return{c(){t=d("h3"),e=R("Auth with OTP ("),h=R(l),i=R(")"),c=T(),n=d("div"),n.innerHTML=`<p>Authenticate with an one-time password (OTP).</p> <p>Note that when requesting an OTP we return an <code>otpId</code> even if a user with the provided email
        doesn&#39;t exist as a very basic enumeration protection.</p>`,m=T(),x(v.$$.fragment),C=T(),F=d("h6"),F.textContent="API details",D=T(),I=d("div"),M=d("div");for(let o=0;o<O.length;o+=1)O[o].c();Q=T(),B=d("div");for(let o=0;o<k.length;o+=1)k[o].c();g(t,"class","m-b-sm"),g(n,"class","content txt-lg m-b-sm"),g(F,"class","m-b-xs"),g(M,"class","tabs-header compact"),g(B,"class","tabs-content"),g(I,"class","tabs")},m(o,$){b(o,t,$),s(t,e),s(t,h),s(t,i),b(o,c,$),b(o,n,$),b(o,m,$),ee(v,o,$),b(o,C,$),b(o,F,$),b(o,D,$),b(o,I,$),s(I,M);for(let P=0;P<O.length;P+=1)O[P]&&O[P].m(M,null);s(I,Q),s(I,B);for(let P=0;P<k.length;P+=1)k[P]&&k[P].m(B,null);H=!0},p(o,[$]){var X,u,p,S;(!H||$&1)&&l!==(l=o[0].name+"")&&ce(h,l);const P={};if($&5&&(P.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${o[2]}');

        ...

        // send OTP email to the provided auth record
        const req = await pb.collection('${(X=o[0])==null?void 0:X.name}').requestOTP('test@example.com');

        // ... show a screen/popup to enter the password from the email ...

        // authenticate with the requested OTP id and the email password
        const authData = await pb.collection('${(u=o[0])==null?void 0:u.name}').authWithOTP(
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

        final pb = PocketBase('${o[2]}');

        ...

        // send OTP email to the provided auth record
        final req = await pb.collection('${(p=o[0])==null?void 0:p.name}').requestOTP('test@example.com');

        // ... show a screen/popup to enter the password from the email ...

        // authenticate with the requested OTP id and the email password
        final authData = await pb.collection('${(S=o[0])==null?void 0:S.name}').authWithOTP(
            req.otpId,
            "YOUR_OTP",
        );

        // after the above you can also access the auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.record.id);

        // "logout"
        pb.authStore.clear();
    `),v.$set(P),$&10){y=V(o[3]);let w;for(w=0;w<y.length;w+=1){const K=Ee(o,y,w);O[w]?O[w].p(K,$):(O[w]=Je(K),O[w].c(),O[w].m(M,null))}for(;w<O.length;w+=1)O[w].d(1);O.length=y.length}if($&11){q=V(o[3]);let w;for(w=0;w<q.length;w+=1){const K=Ve(o,q,w);k[w]?(k[w].p(K,$),E(k[w],1)):(k[w]=Ne(K),k[w].c(),E(k[w],1),k[w].m(B,null))}for(ue(),w=q.length;w<k.length;w+=1)L(w);he()}},i(o){if(!H){E(v.$$.fragment,o);for(let $=0;$<q.length;$+=1)E(k[$]);H=!0}},o(o){j(v.$$.fragment,o),k=k.filter(Boolean);for(let $=0;$<k.length;$+=1)j(k[$]);H=!1},d(o){o&&(_(t),_(c),_(n),_(m),_(C),_(F),_(D),_(I)),te(v,o),Me(O,o),Me(k,o)}}}function at(a,t,e){let l,{collection:h}=t;const i=[{title:"OTP Request",component:lt},{title:"OTP Auth",component:xe}];let c=0;const n=m=>e(1,c=m);return a.$$set=m=>{"collection"in m&&e(0,h=m.collection)},e(2,l=ke.getApiExampleUrl(Ge.baseURL)),[h,c,l,i,n]}class it extends be{constructor(t){super(),_e(this,t,at,ot,ve,{collection:0})}}export{it as default};
