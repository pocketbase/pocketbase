import{S as ve,i as ye,s as $e,M as we,N as nt,e as o,w as f,b as d,c as ot,f as m,g as r,h as e,m as st,x as Dt,P as fe,Q as Pe,k as Re,R as Ce,n as Ae,t as Z,a as x,o as c,d as it,U as Oe,C as pe,p as Te,r as rt,u as Ue}from"./index-a084d9d7.js";import{S as Me}from"./SdkTabs-ba0ec979.js";import{F as De}from"./FieldsQueryParam-71e01e64.js";function he(s,l,a){const i=s.slice();return i[8]=l[a],i}function be(s,l,a){const i=s.slice();return i[8]=l[a],i}function Ee(s){let l;return{c(){l=f("email")},m(a,i){r(a,l,i)},d(a){a&&c(l)}}}function We(s){let l;return{c(){l=f("username")},m(a,i){r(a,l,i)},d(a){a&&c(l)}}}function Le(s){let l;return{c(){l=f("username/email")},m(a,i){r(a,l,i)},d(a){a&&c(l)}}}function me(s){let l;return{c(){l=o("strong"),l.textContent="username"},m(a,i){r(a,l,i)},d(a){a&&c(l)}}}function _e(s){let l;return{c(){l=f("or")},m(a,i){r(a,l,i)},d(a){a&&c(l)}}}function ke(s){let l;return{c(){l=o("strong"),l.textContent="email"},m(a,i){r(a,l,i)},d(a){a&&c(l)}}}function ge(s,l){let a,i=l[8].code+"",g,b,p,u;function _(){return l[7](l[8])}return{key:s,first:null,c(){a=o("button"),g=f(i),b=d(),m(a,"class","tab-item"),rt(a,"active",l[3]===l[8].code),this.first=a},m(R,C){r(R,a,C),e(a,g),e(a,b),p||(u=Ue(a,"click",_),p=!0)},p(R,C){l=R,C&16&&i!==(i=l[8].code+"")&&Dt(g,i),C&24&&rt(a,"active",l[3]===l[8].code)},d(R){R&&c(a),p=!1,u()}}}function Se(s,l){let a,i,g,b;return i=new we({props:{content:l[8].body}}),{key:s,first:null,c(){a=o("div"),ot(i.$$.fragment),g=d(),m(a,"class","tab-item"),rt(a,"active",l[3]===l[8].code),this.first=a},m(p,u){r(p,a,u),st(i,a,null),e(a,g),b=!0},p(p,u){l=p;const _={};u&16&&(_.content=l[8].body),i.$set(_),(!b||u&24)&&rt(a,"active",l[3]===l[8].code)},i(p){b||(Z(i.$$.fragment,p),b=!0)},o(p){x(i.$$.fragment,p),b=!1},d(p){p&&c(a),it(i)}}}function Be(s){var re,ce;let l,a,i=s[0].name+"",g,b,p,u,_,R,C,A,B,Et,ct,T,dt,N,ut,U,tt,Wt,et,I,Lt,ft,lt=s[0].name+"",pt,Bt,ht,V,bt,M,mt,qt,Q,D,_t,Ft,kt,Ht,$,Yt,gt,St,wt,Nt,vt,yt,j,$t,E,Pt,It,J,W,Rt,Vt,Ct,Qt,k,jt,q,Jt,Kt,zt,At,Gt,Ot,Xt,Zt,xt,Tt,te,ee,F,Ut,K,Mt,L,z,O=[],le=new Map,ae,G,S=[],ne=new Map,H;function oe(t,n){if(t[1]&&t[2])return Le;if(t[1])return We;if(t[2])return Ee}let Y=oe(s),P=Y&&Y(s);T=new Me({props:{js:`
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
    `}});let w=s[1]&&me(),v=s[1]&&s[2]&&_e(),y=s[2]&&ke();q=new we({props:{content:"?expand=relField1,relField2.subRelField"}}),F=new De({});let at=nt(s[4]);const se=t=>t[8].code;for(let t=0;t<at.length;t+=1){let n=be(s,at,t),h=se(n);le.set(h,O[t]=ge(h,n))}let X=nt(s[4]);const ie=t=>t[8].code;for(let t=0;t<X.length;t+=1){let n=he(s,X,t),h=ie(n);ne.set(h,S[t]=Se(h,n))}return{c(){l=o("h3"),a=f("Auth with password ("),g=f(i),b=f(")"),p=d(),u=o("div"),_=o("p"),R=f(`Returns new auth token and account data by a combination of
        `),C=o("strong"),P&&P.c(),A=f(`
        and `),B=o("strong"),B.textContent="password",Et=f("."),ct=d(),ot(T.$$.fragment),dt=d(),N=o("h6"),N.textContent="API details",ut=d(),U=o("div"),tt=o("strong"),tt.textContent="POST",Wt=d(),et=o("div"),I=o("p"),Lt=f("/api/collections/"),ft=o("strong"),pt=f(lt),Bt=f("/auth-with-password"),ht=d(),V=o("div"),V.textContent="Body Parameters",bt=d(),M=o("table"),mt=o("thead"),mt.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr>',qt=d(),Q=o("tbody"),D=o("tr"),_t=o("td"),_t.innerHTML='<div class="inline-flex"><span class="label label-success">Required</span> <span>identity</span></div>',Ft=d(),kt=o("td"),kt.innerHTML='<span class="label">String</span>',Ht=d(),$=o("td"),Yt=f(`The
                `),w&&w.c(),gt=d(),v&&v.c(),St=d(),y&&y.c(),wt=f(`
                of the record to authenticate.`),Nt=d(),vt=o("tr"),vt.innerHTML='<td><div class="inline-flex"><span class="label label-success">Required</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>The auth record password.</td>',yt=d(),j=o("div"),j.textContent="Query parameters",$t=d(),E=o("table"),Pt=o("thead"),Pt.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',It=d(),J=o("tbody"),W=o("tr"),Rt=o("td"),Rt.textContent="expand",Vt=d(),Ct=o("td"),Ct.innerHTML='<span class="label">String</span>',Qt=d(),k=o("td"),jt=f(`Auto expand record relations. Ex.:
                `),ot(q.$$.fragment),Jt=f(`
                Supports up to 6-levels depth nested relations expansion. `),Kt=o("br"),zt=f(`
                The expanded relations will be appended to the record under the
                `),At=o("code"),At.textContent="expand",Gt=f(" property (eg. "),Ot=o("code"),Ot.textContent='"expand": {"relField1": {...}, ...}',Xt=f(`).
                `),Zt=o("br"),xt=f(`
                Only the relations to which the request user has permissions to `),Tt=o("strong"),Tt.textContent="view",te=f(" will be expanded."),ee=d(),ot(F.$$.fragment),Ut=d(),K=o("div"),K.textContent="Responses",Mt=d(),L=o("div"),z=o("div");for(let t=0;t<O.length;t+=1)O[t].c();ae=d(),G=o("div");for(let t=0;t<S.length;t+=1)S[t].c();m(l,"class","m-b-sm"),m(u,"class","content txt-lg m-b-sm"),m(N,"class","m-b-xs"),m(tt,"class","label label-primary"),m(et,"class","content"),m(U,"class","alert alert-success"),m(V,"class","section-title"),m(M,"class","table-compact table-border m-b-base"),m(j,"class","section-title"),m(E,"class","table-compact table-border m-b-base"),m(K,"class","section-title"),m(z,"class","tabs-header compact left"),m(G,"class","tabs-content"),m(L,"class","tabs")},m(t,n){r(t,l,n),e(l,a),e(l,g),e(l,b),r(t,p,n),r(t,u,n),e(u,_),e(_,R),e(_,C),P&&P.m(C,null),e(_,A),e(_,B),e(_,Et),r(t,ct,n),st(T,t,n),r(t,dt,n),r(t,N,n),r(t,ut,n),r(t,U,n),e(U,tt),e(U,Wt),e(U,et),e(et,I),e(I,Lt),e(I,ft),e(ft,pt),e(I,Bt),r(t,ht,n),r(t,V,n),r(t,bt,n),r(t,M,n),e(M,mt),e(M,qt),e(M,Q),e(Q,D),e(D,_t),e(D,Ft),e(D,kt),e(D,Ht),e(D,$),e($,Yt),w&&w.m($,null),e($,gt),v&&v.m($,null),e($,St),y&&y.m($,null),e($,wt),e(Q,Nt),e(Q,vt),r(t,yt,n),r(t,j,n),r(t,$t,n),r(t,E,n),e(E,Pt),e(E,It),e(E,J),e(J,W),e(W,Rt),e(W,Vt),e(W,Ct),e(W,Qt),e(W,k),e(k,jt),st(q,k,null),e(k,Jt),e(k,Kt),e(k,zt),e(k,At),e(k,Gt),e(k,Ot),e(k,Xt),e(k,Zt),e(k,xt),e(k,Tt),e(k,te),e(J,ee),st(F,J,null),r(t,Ut,n),r(t,K,n),r(t,Mt,n),r(t,L,n),e(L,z);for(let h=0;h<O.length;h+=1)O[h]&&O[h].m(z,null);e(L,ae),e(L,G);for(let h=0;h<S.length;h+=1)S[h]&&S[h].m(G,null);H=!0},p(t,[n]){var de,ue;(!H||n&1)&&i!==(i=t[0].name+"")&&Dt(g,i),Y!==(Y=oe(t))&&(P&&P.d(1),P=Y&&Y(t),P&&(P.c(),P.m(C,null)));const h={};n&97&&(h.js=`
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
    `),n&97&&(h.dart=`
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
    `),T.$set(h),(!H||n&1)&&lt!==(lt=t[0].name+"")&&Dt(pt,lt),t[1]?w||(w=me(),w.c(),w.m($,gt)):w&&(w.d(1),w=null),t[1]&&t[2]?v||(v=_e(),v.c(),v.m($,St)):v&&(v.d(1),v=null),t[2]?y||(y=ke(),y.c(),y.m($,wt)):y&&(y.d(1),y=null),n&24&&(at=nt(t[4]),O=fe(O,n,se,1,t,at,le,z,Pe,ge,null,be)),n&24&&(X=nt(t[4]),Re(),S=fe(S,n,ie,1,t,X,ne,G,Ce,Se,null,he),Ae())},i(t){if(!H){Z(T.$$.fragment,t),Z(q.$$.fragment,t),Z(F.$$.fragment,t);for(let n=0;n<X.length;n+=1)Z(S[n]);H=!0}},o(t){x(T.$$.fragment,t),x(q.$$.fragment,t),x(F.$$.fragment,t);for(let n=0;n<S.length;n+=1)x(S[n]);H=!1},d(t){t&&(c(l),c(p),c(u),c(ct),c(dt),c(N),c(ut),c(U),c(ht),c(V),c(bt),c(M),c(yt),c(j),c($t),c(E),c(Ut),c(K),c(Mt),c(L)),P&&P.d(),it(T,t),w&&w.d(),v&&v.d(),y&&y.d(),it(q),it(F);for(let n=0;n<O.length;n+=1)O[n].d();for(let n=0;n<S.length;n+=1)S[n].d()}}}function qe(s,l,a){let i,g,b,p,{collection:u=new Oe}=l,_=200,R=[];const C=A=>a(3,_=A.code);return s.$$set=A=>{"collection"in A&&a(0,u=A.collection)},s.$$.update=()=>{var A,B;s.$$.dirty&1&&a(2,g=(A=u==null?void 0:u.options)==null?void 0:A.allowEmailAuth),s.$$.dirty&1&&a(1,b=(B=u==null?void 0:u.options)==null?void 0:B.allowUsernameAuth),s.$$.dirty&6&&a(5,p=b&&g?"YOUR_USERNAME_OR_EMAIL":b?"YOUR_USERNAME":"YOUR_EMAIL"),s.$$.dirty&1&&a(4,R=[{code:200,body:JSON.stringify({token:"JWT_TOKEN",record:pe.dummyCollectionRecord(u)},null,2)},{code:400,body:`
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
            `}])},a(6,i=pe.getApiExampleUrl(Te.baseUrl)),[u,b,g,_,R,p,i,C]}class Ne extends ve{constructor(l){super(),ye(this,l,qe,Be,$e,{collection:0})}}export{Ne as default};
