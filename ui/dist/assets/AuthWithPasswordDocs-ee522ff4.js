import{S as we,i as ye,s as $e,N as ve,O as ot,e as n,w as p,b as d,c as nt,f as m,g as r,h as e,m as st,x as Dt,P as pe,Q as Pe,k as Re,R as Ce,n as Oe,t as Z,a as x,o as c,d as it,C as fe,p as Ae,r as rt,u as Te}from"./index-8354bde7.js";import{S as Ue}from"./SdkTabs-86785e52.js";import{F as Me}from"./FieldsQueryParam-7cb62521.js";function he(s,l,a){const i=s.slice();return i[8]=l[a],i}function be(s,l,a){const i=s.slice();return i[8]=l[a],i}function De(s){let l;return{c(){l=p("email")},m(a,i){r(a,l,i)},d(a){a&&c(l)}}}function Ee(s){let l;return{c(){l=p("username")},m(a,i){r(a,l,i)},d(a){a&&c(l)}}}function We(s){let l;return{c(){l=p("username/email")},m(a,i){r(a,l,i)},d(a){a&&c(l)}}}function me(s){let l;return{c(){l=n("strong"),l.textContent="username"},m(a,i){r(a,l,i)},d(a){a&&c(l)}}}function _e(s){let l;return{c(){l=p("or")},m(a,i){r(a,l,i)},d(a){a&&c(l)}}}function ke(s){let l;return{c(){l=n("strong"),l.textContent="email"},m(a,i){r(a,l,i)},d(a){a&&c(l)}}}function ge(s,l){let a,i=l[8].code+"",g,b,f,u;function _(){return l[7](l[8])}return{key:s,first:null,c(){a=n("button"),g=p(i),b=d(),m(a,"class","tab-item"),rt(a,"active",l[3]===l[8].code),this.first=a},m(R,C){r(R,a,C),e(a,g),e(a,b),f||(u=Te(a,"click",_),f=!0)},p(R,C){l=R,C&16&&i!==(i=l[8].code+"")&&Dt(g,i),C&24&&rt(a,"active",l[3]===l[8].code)},d(R){R&&c(a),f=!1,u()}}}function Se(s,l){let a,i,g,b;return i=new ve({props:{content:l[8].body}}),{key:s,first:null,c(){a=n("div"),nt(i.$$.fragment),g=d(),m(a,"class","tab-item"),rt(a,"active",l[3]===l[8].code),this.first=a},m(f,u){r(f,a,u),st(i,a,null),e(a,g),b=!0},p(f,u){l=f;const _={};u&16&&(_.content=l[8].body),i.$set(_),(!b||u&24)&&rt(a,"active",l[3]===l[8].code)},i(f){b||(Z(i.$$.fragment,f),b=!0)},o(f){x(i.$$.fragment,f),b=!1},d(f){f&&c(a),it(i)}}}function Le(s){var re,ce;let l,a,i=s[0].name+"",g,b,f,u,_,R,C,O,B,Et,ct,T,dt,N,ut,U,tt,Wt,et,I,Lt,pt,lt=s[0].name+"",ft,Bt,ht,V,bt,M,mt,qt,Q,D,_t,Ft,kt,Ht,$,Yt,gt,St,vt,Nt,wt,yt,j,$t,E,Pt,It,J,W,Rt,Vt,Ct,Qt,k,jt,q,Jt,Kt,zt,Ot,Gt,At,Xt,Zt,xt,Tt,te,ee,F,Ut,K,Mt,L,z,A=[],le=new Map,ae,G,S=[],oe=new Map,H;function ne(t,o){if(t[1]&&t[2])return We;if(t[1])return Ee;if(t[2])return De}let Y=ne(s),P=Y&&Y(s);T=new Ue({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${s[6]}');

        ...

        const authData = await pb.collection('${(re=s[0])==null?void 0:re.name}').authWithPassword(
            '${s[5]}',
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

        final pb = PocketBase('${s[6]}');

        ...

        final authData = await pb.collection('${(ce=s[0])==null?void 0:ce.name}').authWithPassword(
          '${s[5]}',
          'YOUR_PASSWORD',
        );

        // after the above you can also access the auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.model.id);

        // "logout" the last authenticated account
        pb.authStore.clear();
    `}});let v=s[1]&&me(),w=s[1]&&s[2]&&_e(),y=s[2]&&ke();q=new ve({props:{content:"?expand=relField1,relField2.subRelField"}}),F=new Me({});let at=ot(s[4]);const se=t=>t[8].code;for(let t=0;t<at.length;t+=1){let o=be(s,at,t),h=se(o);le.set(h,A[t]=ge(h,o))}let X=ot(s[4]);const ie=t=>t[8].code;for(let t=0;t<X.length;t+=1){let o=he(s,X,t),h=ie(o);oe.set(h,S[t]=Se(h,o))}return{c(){l=n("h3"),a=p("Auth with password ("),g=p(i),b=p(")"),f=d(),u=n("div"),_=n("p"),R=p(`Returns new auth token and account data by a combination of
        `),C=n("strong"),P&&P.c(),O=p(`
        and `),B=n("strong"),B.textContent="password",Et=p("."),ct=d(),nt(T.$$.fragment),dt=d(),N=n("h6"),N.textContent="API details",ut=d(),U=n("div"),tt=n("strong"),tt.textContent="POST",Wt=d(),et=n("div"),I=n("p"),Lt=p("/api/collections/"),pt=n("strong"),ft=p(lt),Bt=p("/auth-with-password"),ht=d(),V=n("div"),V.textContent="Body Parameters",bt=d(),M=n("table"),mt=n("thead"),mt.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr>',qt=d(),Q=n("tbody"),D=n("tr"),_t=n("td"),_t.innerHTML='<div class="inline-flex"><span class="label label-success">Required</span> <span>identity</span></div>',Ft=d(),kt=n("td"),kt.innerHTML='<span class="label">String</span>',Ht=d(),$=n("td"),Yt=p(`The
                `),v&&v.c(),gt=d(),w&&w.c(),St=d(),y&&y.c(),vt=p(`
                of the record to authenticate.`),Nt=d(),wt=n("tr"),wt.innerHTML='<td><div class="inline-flex"><span class="label label-success">Required</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>The auth record password.</td>',yt=d(),j=n("div"),j.textContent="Query parameters",$t=d(),E=n("table"),Pt=n("thead"),Pt.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',It=d(),J=n("tbody"),W=n("tr"),Rt=n("td"),Rt.textContent="expand",Vt=d(),Ct=n("td"),Ct.innerHTML='<span class="label">String</span>',Qt=d(),k=n("td"),jt=p(`Auto expand record relations. Ex.:
                `),nt(q.$$.fragment),Jt=p(`
                Supports up to 6-levels depth nested relations expansion. `),Kt=n("br"),zt=p(`
                The expanded relations will be appended to the record under the
                `),Ot=n("code"),Ot.textContent="expand",Gt=p(" property (eg. "),At=n("code"),At.textContent='"expand": {"relField1": {...}, ...}',Xt=p(`).
                `),Zt=n("br"),xt=p(`
                Only the relations to which the request user has permissions to `),Tt=n("strong"),Tt.textContent="view",te=p(" will be expanded."),ee=d(),nt(F.$$.fragment),Ut=d(),K=n("div"),K.textContent="Responses",Mt=d(),L=n("div"),z=n("div");for(let t=0;t<A.length;t+=1)A[t].c();ae=d(),G=n("div");for(let t=0;t<S.length;t+=1)S[t].c();m(l,"class","m-b-sm"),m(u,"class","content txt-lg m-b-sm"),m(N,"class","m-b-xs"),m(tt,"class","label label-primary"),m(et,"class","content"),m(U,"class","alert alert-success"),m(V,"class","section-title"),m(M,"class","table-compact table-border m-b-base"),m(j,"class","section-title"),m(E,"class","table-compact table-border m-b-base"),m(K,"class","section-title"),m(z,"class","tabs-header compact combined left"),m(G,"class","tabs-content"),m(L,"class","tabs")},m(t,o){r(t,l,o),e(l,a),e(l,g),e(l,b),r(t,f,o),r(t,u,o),e(u,_),e(_,R),e(_,C),P&&P.m(C,null),e(_,O),e(_,B),e(_,Et),r(t,ct,o),st(T,t,o),r(t,dt,o),r(t,N,o),r(t,ut,o),r(t,U,o),e(U,tt),e(U,Wt),e(U,et),e(et,I),e(I,Lt),e(I,pt),e(pt,ft),e(I,Bt),r(t,ht,o),r(t,V,o),r(t,bt,o),r(t,M,o),e(M,mt),e(M,qt),e(M,Q),e(Q,D),e(D,_t),e(D,Ft),e(D,kt),e(D,Ht),e(D,$),e($,Yt),v&&v.m($,null),e($,gt),w&&w.m($,null),e($,St),y&&y.m($,null),e($,vt),e(Q,Nt),e(Q,wt),r(t,yt,o),r(t,j,o),r(t,$t,o),r(t,E,o),e(E,Pt),e(E,It),e(E,J),e(J,W),e(W,Rt),e(W,Vt),e(W,Ct),e(W,Qt),e(W,k),e(k,jt),st(q,k,null),e(k,Jt),e(k,Kt),e(k,zt),e(k,Ot),e(k,Gt),e(k,At),e(k,Xt),e(k,Zt),e(k,xt),e(k,Tt),e(k,te),e(J,ee),st(F,J,null),r(t,Ut,o),r(t,K,o),r(t,Mt,o),r(t,L,o),e(L,z);for(let h=0;h<A.length;h+=1)A[h]&&A[h].m(z,null);e(L,ae),e(L,G);for(let h=0;h<S.length;h+=1)S[h]&&S[h].m(G,null);H=!0},p(t,[o]){var de,ue;(!H||o&1)&&i!==(i=t[0].name+"")&&Dt(g,i),Y!==(Y=ne(t))&&(P&&P.d(1),P=Y&&Y(t),P&&(P.c(),P.m(C,null)));const h={};o&97&&(h.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${t[6]}');

        ...

        const authData = await pb.collection('${(de=t[0])==null?void 0:de.name}').authWithPassword(
            '${t[5]}',
            'YOUR_PASSWORD',
        );

        // after the above you can also access the auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);

        // "logout" the last authenticated account
        pb.authStore.clear();
    `),o&97&&(h.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${t[6]}');

        ...

        final authData = await pb.collection('${(ue=t[0])==null?void 0:ue.name}').authWithPassword(
          '${t[5]}',
          'YOUR_PASSWORD',
        );

        // after the above you can also access the auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.model.id);

        // "logout" the last authenticated account
        pb.authStore.clear();
    `),T.$set(h),(!H||o&1)&&lt!==(lt=t[0].name+"")&&Dt(ft,lt),t[1]?v||(v=me(),v.c(),v.m($,gt)):v&&(v.d(1),v=null),t[1]&&t[2]?w||(w=_e(),w.c(),w.m($,St)):w&&(w.d(1),w=null),t[2]?y||(y=ke(),y.c(),y.m($,vt)):y&&(y.d(1),y=null),o&24&&(at=ot(t[4]),A=pe(A,o,se,1,t,at,le,z,Pe,ge,null,be)),o&24&&(X=ot(t[4]),Re(),S=pe(S,o,ie,1,t,X,oe,G,Ce,Se,null,he),Oe())},i(t){if(!H){Z(T.$$.fragment,t),Z(q.$$.fragment,t),Z(F.$$.fragment,t);for(let o=0;o<X.length;o+=1)Z(S[o]);H=!0}},o(t){x(T.$$.fragment,t),x(q.$$.fragment,t),x(F.$$.fragment,t);for(let o=0;o<S.length;o+=1)x(S[o]);H=!1},d(t){t&&(c(l),c(f),c(u),c(ct),c(dt),c(N),c(ut),c(U),c(ht),c(V),c(bt),c(M),c(yt),c(j),c($t),c(E),c(Ut),c(K),c(Mt),c(L)),P&&P.d(),it(T,t),v&&v.d(),w&&w.d(),y&&y.d(),it(q),it(F);for(let o=0;o<A.length;o+=1)A[o].d();for(let o=0;o<S.length;o+=1)S[o].d()}}}function Be(s,l,a){let i,g,b,f,{collection:u}=l,_=200,R=[];const C=O=>a(3,_=O.code);return s.$$set=O=>{"collection"in O&&a(0,u=O.collection)},s.$$.update=()=>{var O,B;s.$$.dirty&1&&a(2,g=(O=u==null?void 0:u.options)==null?void 0:O.allowEmailAuth),s.$$.dirty&1&&a(1,b=(B=u==null?void 0:u.options)==null?void 0:B.allowUsernameAuth),s.$$.dirty&6&&a(5,f=b&&g?"YOUR_USERNAME_OR_EMAIL":b?"YOUR_USERNAME":"YOUR_EMAIL"),s.$$.dirty&1&&a(4,R=[{code:200,body:JSON.stringify({token:"JWT_TOKEN",record:fe.dummyCollectionRecord(u)},null,2)},{code:400,body:`
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
            `}])},a(6,i=fe.getApiExampleUrl(Ae.baseUrl)),[u,b,g,_,R,f,i,C]}class Ye extends we{constructor(l){super(),ye(this,l,Be,Le,$e,{collection:0})}}export{Ye as default};
