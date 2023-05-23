import{S as ve,i as ye,s as ge,M as we,e as s,w as f,b as d,c as ot,f as h,g as r,h as e,m as st,x as Ut,N as ue,P as $e,k as Pe,Q as Re,n as Ce,t as Z,a as x,o as c,d as nt,T as Te,C as fe,p as Ae,r as it,u as Oe}from"./index-a65ca895.js";import{S as Me}from"./SdkTabs-ad912c8f.js";import{F as Ue}from"./FieldsQueryParam-ba250473.js";function pe(n,l,a){const i=n.slice();return i[8]=l[a],i}function be(n,l,a){const i=n.slice();return i[8]=l[a],i}function De(n){let l;return{c(){l=f("email")},m(a,i){r(a,l,i)},d(a){a&&c(l)}}}function Ee(n){let l;return{c(){l=f("username")},m(a,i){r(a,l,i)},d(a){a&&c(l)}}}function We(n){let l;return{c(){l=f("username/email")},m(a,i){r(a,l,i)},d(a){a&&c(l)}}}function me(n){let l;return{c(){l=s("strong"),l.textContent="username"},m(a,i){r(a,l,i)},d(a){a&&c(l)}}}function he(n){let l;return{c(){l=f("or")},m(a,i){r(a,l,i)},d(a){a&&c(l)}}}function _e(n){let l;return{c(){l=s("strong"),l.textContent="email"},m(a,i){r(a,l,i)},d(a){a&&c(l)}}}function ke(n,l){let a,i=l[8].code+"",S,m,p,u;function _(){return l[7](l[8])}return{key:n,first:null,c(){a=s("button"),S=f(i),m=d(),h(a,"class","tab-item"),it(a,"active",l[3]===l[8].code),this.first=a},m(R,C){r(R,a,C),e(a,S),e(a,m),p||(u=Oe(a,"click",_),p=!0)},p(R,C){l=R,C&16&&i!==(i=l[8].code+"")&&Ut(S,i),C&24&&it(a,"active",l[3]===l[8].code)},d(R){R&&c(a),p=!1,u()}}}function Se(n,l){let a,i,S,m;return i=new we({props:{content:l[8].body}}),{key:n,first:null,c(){a=s("div"),ot(i.$$.fragment),S=d(),h(a,"class","tab-item"),it(a,"active",l[3]===l[8].code),this.first=a},m(p,u){r(p,a,u),st(i,a,null),e(a,S),m=!0},p(p,u){l=p;const _={};u&16&&(_.content=l[8].body),i.$set(_),(!m||u&24)&&it(a,"active",l[3]===l[8].code)},i(p){m||(Z(i.$$.fragment,p),m=!0)},o(p){x(i.$$.fragment,p),m=!1},d(p){p&&c(a),nt(i)}}}function Le(n){var ie,re;let l,a,i=n[0].name+"",S,m,p,u,_,R,C,T,B,Dt,rt,O,ct,N,dt,M,tt,Et,et,I,Wt,ut,lt=n[0].name+"",ft,Lt,pt,V,bt,U,mt,Bt,Q,D,ht,qt,_t,Ft,$,Ht,kt,St,wt,Yt,vt,yt,j,gt,E,$t,Nt,J,W,Pt,It,Rt,Vt,k,Qt,q,jt,Jt,Kt,Ct,zt,Tt,Gt,Xt,Zt,At,xt,te,F,Ot,K,Mt,L,z,A=[],ee=new Map,le,G,w=[],ae=new Map,H;function oe(t,o){if(t[1]&&t[2])return We;if(t[1])return Ee;if(t[2])return De}let Y=oe(n),P=Y&&Y(n);O=new Me({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${n[6]}');

        ...

        const authData = await pb.collection('${(ie=n[0])==null?void 0:ie.name}').authWithPassword(
            '${n[5]}',
            'YOUR_PASSWORD',
        );

        // after the above you can also access the auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);

        // "logout" the last authenticated account
        pb.authStore.clear();
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${n[6]}');

        ...

        final authData = await pb.collection('${(re=n[0])==null?void 0:re.name}').authWithPassword(
          '${n[5]}',
          'YOUR_PASSWORD',
        );

        // after the above you can also access the auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.model.id);

        // "logout" the last authenticated account
        pb.authStore.clear();
    `}});let v=n[1]&&me(),y=n[1]&&n[2]&&he(),g=n[2]&&_e();q=new we({props:{content:"?expand=relField1,relField2.subRelField"}}),F=new Ue({});let at=n[4];const se=t=>t[8].code;for(let t=0;t<at.length;t+=1){let o=be(n,at,t),b=se(o);ee.set(b,A[t]=ke(b,o))}let X=n[4];const ne=t=>t[8].code;for(let t=0;t<X.length;t+=1){let o=pe(n,X,t),b=ne(o);ae.set(b,w[t]=Se(b,o))}return{c(){l=s("h3"),a=f("Auth with password ("),S=f(i),m=f(")"),p=d(),u=s("div"),_=s("p"),R=f(`Returns new auth token and account data by a combination of
        `),C=s("strong"),P&&P.c(),T=f(`
        and `),B=s("strong"),B.textContent="password",Dt=f("."),rt=d(),ot(O.$$.fragment),ct=d(),N=s("h6"),N.textContent="API details",dt=d(),M=s("div"),tt=s("strong"),tt.textContent="POST",Et=d(),et=s("div"),I=s("p"),Wt=f("/api/collections/"),ut=s("strong"),ft=f(lt),Lt=f("/auth-with-password"),pt=d(),V=s("div"),V.textContent="Body Parameters",bt=d(),U=s("table"),mt=s("thead"),mt.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr>`,Bt=d(),Q=s("tbody"),D=s("tr"),ht=s("td"),ht.innerHTML=`<div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>identity</span></div>`,qt=d(),_t=s("td"),_t.innerHTML='<span class="label">String</span>',Ft=d(),$=s("td"),Ht=f(`The
                `),v&&v.c(),kt=d(),y&&y.c(),St=d(),g&&g.c(),wt=f(`
                of the record to authenticate.`),Yt=d(),vt=s("tr"),vt.innerHTML=`<td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>password</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The auth record password.</td>`,yt=d(),j=s("div"),j.textContent="Query parameters",gt=d(),E=s("table"),$t=s("thead"),$t.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,Nt=d(),J=s("tbody"),W=s("tr"),Pt=s("td"),Pt.textContent="expand",It=d(),Rt=s("td"),Rt.innerHTML='<span class="label">String</span>',Vt=d(),k=s("td"),Qt=f(`Auto expand record relations. Ex.:
                `),ot(q.$$.fragment),jt=f(`
                Supports up to 6-levels depth nested relations expansion. `),Jt=s("br"),Kt=f(`
                The expanded relations will be appended to the record under the
                `),Ct=s("code"),Ct.textContent="expand",zt=f(" property (eg. "),Tt=s("code"),Tt.textContent='"expand": {"relField1": {...}, ...}',Gt=f(`).
                `),Xt=s("br"),Zt=f(`
                Only the relations to which the request user has permissions to `),At=s("strong"),At.textContent="view",xt=f(" will be expanded."),te=d(),ot(F.$$.fragment),Ot=d(),K=s("div"),K.textContent="Responses",Mt=d(),L=s("div"),z=s("div");for(let t=0;t<A.length;t+=1)A[t].c();le=d(),G=s("div");for(let t=0;t<w.length;t+=1)w[t].c();h(l,"class","m-b-sm"),h(u,"class","content txt-lg m-b-sm"),h(N,"class","m-b-xs"),h(tt,"class","label label-primary"),h(et,"class","content"),h(M,"class","alert alert-success"),h(V,"class","section-title"),h(U,"class","table-compact table-border m-b-base"),h(j,"class","section-title"),h(E,"class","table-compact table-border m-b-base"),h(K,"class","section-title"),h(z,"class","tabs-header compact left"),h(G,"class","tabs-content"),h(L,"class","tabs")},m(t,o){r(t,l,o),e(l,a),e(l,S),e(l,m),r(t,p,o),r(t,u,o),e(u,_),e(_,R),e(_,C),P&&P.m(C,null),e(_,T),e(_,B),e(_,Dt),r(t,rt,o),st(O,t,o),r(t,ct,o),r(t,N,o),r(t,dt,o),r(t,M,o),e(M,tt),e(M,Et),e(M,et),e(et,I),e(I,Wt),e(I,ut),e(ut,ft),e(I,Lt),r(t,pt,o),r(t,V,o),r(t,bt,o),r(t,U,o),e(U,mt),e(U,Bt),e(U,Q),e(Q,D),e(D,ht),e(D,qt),e(D,_t),e(D,Ft),e(D,$),e($,Ht),v&&v.m($,null),e($,kt),y&&y.m($,null),e($,St),g&&g.m($,null),e($,wt),e(Q,Yt),e(Q,vt),r(t,yt,o),r(t,j,o),r(t,gt,o),r(t,E,o),e(E,$t),e(E,Nt),e(E,J),e(J,W),e(W,Pt),e(W,It),e(W,Rt),e(W,Vt),e(W,k),e(k,Qt),st(q,k,null),e(k,jt),e(k,Jt),e(k,Kt),e(k,Ct),e(k,zt),e(k,Tt),e(k,Gt),e(k,Xt),e(k,Zt),e(k,At),e(k,xt),e(J,te),st(F,J,null),r(t,Ot,o),r(t,K,o),r(t,Mt,o),r(t,L,o),e(L,z);for(let b=0;b<A.length;b+=1)A[b]&&A[b].m(z,null);e(L,le),e(L,G);for(let b=0;b<w.length;b+=1)w[b]&&w[b].m(G,null);H=!0},p(t,[o]){var ce,de;(!H||o&1)&&i!==(i=t[0].name+"")&&Ut(S,i),Y!==(Y=oe(t))&&(P&&P.d(1),P=Y&&Y(t),P&&(P.c(),P.m(C,null)));const b={};o&97&&(b.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${t[6]}');

        ...

        const authData = await pb.collection('${(ce=t[0])==null?void 0:ce.name}').authWithPassword(
            '${t[5]}',
            'YOUR_PASSWORD',
        );

        // after the above you can also access the auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);

        // "logout" the last authenticated account
        pb.authStore.clear();
    `),o&97&&(b.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${t[6]}');

        ...

        final authData = await pb.collection('${(de=t[0])==null?void 0:de.name}').authWithPassword(
          '${t[5]}',
          'YOUR_PASSWORD',
        );

        // after the above you can also access the auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.model.id);

        // "logout" the last authenticated account
        pb.authStore.clear();
    `),O.$set(b),(!H||o&1)&&lt!==(lt=t[0].name+"")&&Ut(ft,lt),t[1]?v||(v=me(),v.c(),v.m($,kt)):v&&(v.d(1),v=null),t[1]&&t[2]?y||(y=he(),y.c(),y.m($,St)):y&&(y.d(1),y=null),t[2]?g||(g=_e(),g.c(),g.m($,wt)):g&&(g.d(1),g=null),o&24&&(at=t[4],A=ue(A,o,se,1,t,at,ee,z,$e,ke,null,be)),o&24&&(X=t[4],Pe(),w=ue(w,o,ne,1,t,X,ae,G,Re,Se,null,pe),Ce())},i(t){if(!H){Z(O.$$.fragment,t),Z(q.$$.fragment,t),Z(F.$$.fragment,t);for(let o=0;o<X.length;o+=1)Z(w[o]);H=!0}},o(t){x(O.$$.fragment,t),x(q.$$.fragment,t),x(F.$$.fragment,t);for(let o=0;o<w.length;o+=1)x(w[o]);H=!1},d(t){t&&c(l),t&&c(p),t&&c(u),P&&P.d(),t&&c(rt),nt(O,t),t&&c(ct),t&&c(N),t&&c(dt),t&&c(M),t&&c(pt),t&&c(V),t&&c(bt),t&&c(U),v&&v.d(),y&&y.d(),g&&g.d(),t&&c(yt),t&&c(j),t&&c(gt),t&&c(E),nt(q),nt(F),t&&c(Ot),t&&c(K),t&&c(Mt),t&&c(L);for(let o=0;o<A.length;o+=1)A[o].d();for(let o=0;o<w.length;o+=1)w[o].d()}}}function Be(n,l,a){let i,S,m,p,{collection:u=new Te}=l,_=200,R=[];const C=T=>a(3,_=T.code);return n.$$set=T=>{"collection"in T&&a(0,u=T.collection)},n.$$.update=()=>{var T,B;n.$$.dirty&1&&a(2,S=(T=u==null?void 0:u.options)==null?void 0:T.allowEmailAuth),n.$$.dirty&1&&a(1,m=(B=u==null?void 0:u.options)==null?void 0:B.allowUsernameAuth),n.$$.dirty&6&&a(5,p=m&&S?"YOUR_USERNAME_OR_EMAIL":m?"YOUR_USERNAME":"YOUR_EMAIL"),n.$$.dirty&1&&a(4,R=[{code:200,body:JSON.stringify({token:"JWT_TOKEN",record:fe.dummyCollectionRecord(u)},null,2)},{code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to authenticate.",
                  "data": {
                    "identity": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `}])},a(6,i=fe.getApiExampleUrl(Ae.baseUrl)),[u,m,S,_,R,p,i,C]}class Ye extends ve{constructor(l){super(),ye(this,l,Be,Le,ge,{collection:0})}}export{Ye as default};
