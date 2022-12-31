import{S as Se,i as ve,s as we,N as ke,e as s,w as f,b as u,c as Ot,f as h,g as r,h as o,m as At,x as Tt,O as ce,P as ye,k as ge,Q as Pe,n as Re,t as tt,a as et,o as c,d as Ut,R as $e,C as de,p as Ce,r as lt,u as Oe}from"./index.89a3f554.js";import{S as Ae}from"./SdkTabs.0a6ad1c9.js";function ue(n,e,l){const i=n.slice();return i[8]=e[l],i}function fe(n,e,l){const i=n.slice();return i[8]=e[l],i}function Te(n){let e;return{c(){e=f("email")},m(l,i){r(l,e,i)},d(l){l&&c(e)}}}function Ue(n){let e;return{c(){e=f("username")},m(l,i){r(l,e,i)},d(l){l&&c(e)}}}function Me(n){let e;return{c(){e=f("username/email")},m(l,i){r(l,e,i)},d(l){l&&c(e)}}}function pe(n){let e;return{c(){e=s("strong"),e.textContent="username"},m(l,i){r(l,e,i)},d(l){l&&c(e)}}}function be(n){let e;return{c(){e=f("or")},m(l,i){r(l,e,i)},d(l){l&&c(e)}}}function me(n){let e;return{c(){e=s("strong"),e.textContent="email"},m(l,i){r(l,e,i)},d(l){l&&c(e)}}}function he(n,e){let l,i=e[8].code+"",S,m,p,d;function _(){return e[7](e[8])}return{key:n,first:null,c(){l=s("button"),S=f(i),m=u(),h(l,"class","tab-item"),lt(l,"active",e[3]===e[8].code),this.first=l},m($,C){r($,l,C),o(l,S),o(l,m),p||(d=Oe(l,"click",_),p=!0)},p($,C){e=$,C&16&&i!==(i=e[8].code+"")&&Tt(S,i),C&24&&lt(l,"active",e[3]===e[8].code)},d($){$&&c(l),p=!1,d()}}}function _e(n,e){let l,i,S,m;return i=new ke({props:{content:e[8].body}}),{key:n,first:null,c(){l=s("div"),Ot(i.$$.fragment),S=u(),h(l,"class","tab-item"),lt(l,"active",e[3]===e[8].code),this.first=l},m(p,d){r(p,l,d),At(i,l,null),o(l,S),m=!0},p(p,d){e=p;const _={};d&16&&(_.content=e[8].body),i.$set(_),(!m||d&24)&&lt(l,"active",e[3]===e[8].code)},i(p){m||(tt(i.$$.fragment,p),m=!0)},o(p){et(i.$$.fragment,p),m=!1},d(p){p&&c(l),Ut(i)}}}function De(n){var se,ne;let e,l,i=n[0].name+"",S,m,p,d,_,$,C,O,B,Mt,ot,T,at,F,st,U,G,Dt,X,I,Et,nt,Z=n[0].name+"",it,Wt,rt,N,ct,M,dt,Lt,V,D,ut,Bt,ft,Ht,g,Yt,pt,bt,mt,qt,ht,_t,j,kt,E,St,Ft,vt,W,wt,It,yt,Nt,k,Vt,H,jt,Jt,Qt,gt,Kt,Pt,zt,Gt,Xt,Rt,Zt,$t,J,Ct,L,Q,A=[],xt=new Map,te,K,P=[],ee=new Map,Y;function le(t,a){if(t[1]&&t[2])return Me;if(t[1])return Ue;if(t[2])return Te}let q=le(n),R=q&&q(n);T=new Ae({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${n[6]}');

        ...

        const authData = await pb.collection('${(se=n[0])==null?void 0:se.name}').authWithPassword(
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

        final authData = await pb.collection('${(ne=n[0])==null?void 0:ne.name}').authWithPassword(
          '${n[5]}',
          'YOUR_PASSWORD',
        );

        // after the above you can also access the auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.model.id);

        // "logout" the last authenticated account
        pb.authStore.clear();
    `}});let v=n[1]&&pe(),w=n[1]&&n[2]&&be(),y=n[2]&&me();H=new ke({props:{content:"?expand=relField1,relField2.subRelField"}});let x=n[4];const oe=t=>t[8].code;for(let t=0;t<x.length;t+=1){let a=fe(n,x,t),b=oe(a);xt.set(b,A[t]=he(b,a))}let z=n[4];const ae=t=>t[8].code;for(let t=0;t<z.length;t+=1){let a=ue(n,z,t),b=ae(a);ee.set(b,P[t]=_e(b,a))}return{c(){e=s("h3"),l=f("Auth with password ("),S=f(i),m=f(")"),p=u(),d=s("div"),_=s("p"),$=f(`Returns new auth token and account data by a combination of
        `),C=s("strong"),R&&R.c(),O=f(`
        and `),B=s("strong"),B.textContent="password",Mt=f("."),ot=u(),Ot(T.$$.fragment),at=u(),F=s("h6"),F.textContent="API details",st=u(),U=s("div"),G=s("strong"),G.textContent="POST",Dt=u(),X=s("div"),I=s("p"),Et=f("/api/collections/"),nt=s("strong"),it=f(Z),Wt=f("/auth-with-password"),rt=u(),N=s("div"),N.textContent="Body Parameters",ct=u(),M=s("table"),dt=s("thead"),dt.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr>`,Lt=u(),V=s("tbody"),D=s("tr"),ut=s("td"),ut.innerHTML=`<div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>identity</span></div>`,Bt=u(),ft=s("td"),ft.innerHTML='<span class="label">String</span>',Ht=u(),g=s("td"),Yt=f(`The
                `),v&&v.c(),pt=u(),w&&w.c(),bt=u(),y&&y.c(),mt=f(`
                of the record to authenticate.`),qt=u(),ht=s("tr"),ht.innerHTML=`<td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>password</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The auth record password.</td>`,_t=u(),j=s("div"),j.textContent="Query parameters",kt=u(),E=s("table"),St=s("thead"),St.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,Ft=u(),vt=s("tbody"),W=s("tr"),wt=s("td"),wt.textContent="expand",It=u(),yt=s("td"),yt.innerHTML='<span class="label">String</span>',Nt=u(),k=s("td"),Vt=f(`Auto expand record relations. Ex.:
                `),Ot(H.$$.fragment),jt=f(`
                Supports up to 6-levels depth nested relations expansion. `),Jt=s("br"),Qt=f(`
                The expanded relations will be appended to the record under the
                `),gt=s("code"),gt.textContent="expand",Kt=f(" property (eg. "),Pt=s("code"),Pt.textContent='"expand": {"relField1": {...}, ...}',zt=f(`).
                `),Gt=s("br"),Xt=f(`
                Only the relations to which the request user has permissions to `),Rt=s("strong"),Rt.textContent="view",Zt=f(" will be expanded."),$t=u(),J=s("div"),J.textContent="Responses",Ct=u(),L=s("div"),Q=s("div");for(let t=0;t<A.length;t+=1)A[t].c();te=u(),K=s("div");for(let t=0;t<P.length;t+=1)P[t].c();h(e,"class","m-b-sm"),h(d,"class","content txt-lg m-b-sm"),h(F,"class","m-b-xs"),h(G,"class","label label-primary"),h(X,"class","content"),h(U,"class","alert alert-success"),h(N,"class","section-title"),h(M,"class","table-compact table-border m-b-base"),h(j,"class","section-title"),h(E,"class","table-compact table-border m-b-base"),h(J,"class","section-title"),h(Q,"class","tabs-header compact left"),h(K,"class","tabs-content"),h(L,"class","tabs")},m(t,a){r(t,e,a),o(e,l),o(e,S),o(e,m),r(t,p,a),r(t,d,a),o(d,_),o(_,$),o(_,C),R&&R.m(C,null),o(_,O),o(_,B),o(_,Mt),r(t,ot,a),At(T,t,a),r(t,at,a),r(t,F,a),r(t,st,a),r(t,U,a),o(U,G),o(U,Dt),o(U,X),o(X,I),o(I,Et),o(I,nt),o(nt,it),o(I,Wt),r(t,rt,a),r(t,N,a),r(t,ct,a),r(t,M,a),o(M,dt),o(M,Lt),o(M,V),o(V,D),o(D,ut),o(D,Bt),o(D,ft),o(D,Ht),o(D,g),o(g,Yt),v&&v.m(g,null),o(g,pt),w&&w.m(g,null),o(g,bt),y&&y.m(g,null),o(g,mt),o(V,qt),o(V,ht),r(t,_t,a),r(t,j,a),r(t,kt,a),r(t,E,a),o(E,St),o(E,Ft),o(E,vt),o(vt,W),o(W,wt),o(W,It),o(W,yt),o(W,Nt),o(W,k),o(k,Vt),At(H,k,null),o(k,jt),o(k,Jt),o(k,Qt),o(k,gt),o(k,Kt),o(k,Pt),o(k,zt),o(k,Gt),o(k,Xt),o(k,Rt),o(k,Zt),r(t,$t,a),r(t,J,a),r(t,Ct,a),r(t,L,a),o(L,Q);for(let b=0;b<A.length;b+=1)A[b].m(Q,null);o(L,te),o(L,K);for(let b=0;b<P.length;b+=1)P[b].m(K,null);Y=!0},p(t,[a]){var ie,re;(!Y||a&1)&&i!==(i=t[0].name+"")&&Tt(S,i),q!==(q=le(t))&&(R&&R.d(1),R=q&&q(t),R&&(R.c(),R.m(C,null)));const b={};a&97&&(b.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${t[6]}');

        ...

        const authData = await pb.collection('${(ie=t[0])==null?void 0:ie.name}').authWithPassword(
            '${t[5]}',
            'YOUR_PASSWORD',
        );

        // after the above you can also access the auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);

        // "logout" the last authenticated account
        pb.authStore.clear();
    `),a&97&&(b.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${t[6]}');

        ...

        final authData = await pb.collection('${(re=t[0])==null?void 0:re.name}').authWithPassword(
          '${t[5]}',
          'YOUR_PASSWORD',
        );

        // after the above you can also access the auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.model.id);

        // "logout" the last authenticated account
        pb.authStore.clear();
    `),T.$set(b),(!Y||a&1)&&Z!==(Z=t[0].name+"")&&Tt(it,Z),t[1]?v||(v=pe(),v.c(),v.m(g,pt)):v&&(v.d(1),v=null),t[1]&&t[2]?w||(w=be(),w.c(),w.m(g,bt)):w&&(w.d(1),w=null),t[2]?y||(y=me(),y.c(),y.m(g,mt)):y&&(y.d(1),y=null),a&24&&(x=t[4],A=ce(A,a,oe,1,t,x,xt,Q,ye,he,null,fe)),a&24&&(z=t[4],ge(),P=ce(P,a,ae,1,t,z,ee,K,Pe,_e,null,ue),Re())},i(t){if(!Y){tt(T.$$.fragment,t),tt(H.$$.fragment,t);for(let a=0;a<z.length;a+=1)tt(P[a]);Y=!0}},o(t){et(T.$$.fragment,t),et(H.$$.fragment,t);for(let a=0;a<P.length;a+=1)et(P[a]);Y=!1},d(t){t&&c(e),t&&c(p),t&&c(d),R&&R.d(),t&&c(ot),Ut(T,t),t&&c(at),t&&c(F),t&&c(st),t&&c(U),t&&c(rt),t&&c(N),t&&c(ct),t&&c(M),v&&v.d(),w&&w.d(),y&&y.d(),t&&c(_t),t&&c(j),t&&c(kt),t&&c(E),Ut(H),t&&c($t),t&&c(J),t&&c(Ct),t&&c(L);for(let a=0;a<A.length;a+=1)A[a].d();for(let a=0;a<P.length;a+=1)P[a].d()}}}function Ee(n,e,l){let i,S,m,p,{collection:d=new $e}=e,_=200,$=[];const C=O=>l(3,_=O.code);return n.$$set=O=>{"collection"in O&&l(0,d=O.collection)},n.$$.update=()=>{var O,B;n.$$.dirty&1&&l(2,S=(O=d==null?void 0:d.options)==null?void 0:O.allowEmailAuth),n.$$.dirty&1&&l(1,m=(B=d==null?void 0:d.options)==null?void 0:B.allowUsernameAuth),n.$$.dirty&6&&l(5,p=m&&S?"YOUR_USERNAME_OR_EMAIL":m?"YOUR_USERNAME":"YOUR_EMAIL"),n.$$.dirty&1&&l(4,$=[{code:200,body:JSON.stringify({token:"JWT_TOKEN",record:de.dummyCollectionRecord(d)},null,2)},{code:400,body:`
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
            `}])},l(6,i=de.getApiExampleUrl(Ce.baseUrl)),[d,m,S,_,$,p,i,C]}class Be extends Se{constructor(e){super(),ve(this,e,Ee,De,we,{collection:0})}}export{Be as default};
