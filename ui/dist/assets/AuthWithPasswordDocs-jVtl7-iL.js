import{S as wt,i as yt,s as $t,N as vt,O as oe,e as n,v as p,b as d,c as ne,f as m,g as r,h as t,m as se,w as De,P as pt,Q as Pt,k as Rt,R as At,n as Ct,t as Z,a as x,o as c,d as ie,C as ft,A as Ot,q as re,r as Tt}from"./index-D0DO79Dq.js";import{S as Ut}from"./SdkTabs-DC6EUYpr.js";import{F as Mt}from"./FieldsQueryParam-BwleQAus.js";function ht(s,l,a){const i=s.slice();return i[8]=l[a],i}function bt(s,l,a){const i=s.slice();return i[8]=l[a],i}function Dt(s){let l;return{c(){l=p("email")},m(a,i){r(a,l,i)},d(a){a&&c(l)}}}function Et(s){let l;return{c(){l=p("username")},m(a,i){r(a,l,i)},d(a){a&&c(l)}}}function Wt(s){let l;return{c(){l=p("username/email")},m(a,i){r(a,l,i)},d(a){a&&c(l)}}}function mt(s){let l;return{c(){l=n("strong"),l.textContent="username"},m(a,i){r(a,l,i)},d(a){a&&c(l)}}}function _t(s){let l;return{c(){l=p("or")},m(a,i){r(a,l,i)},d(a){a&&c(l)}}}function kt(s){let l;return{c(){l=n("strong"),l.textContent="email"},m(a,i){r(a,l,i)},d(a){a&&c(l)}}}function gt(s,l){let a,i=l[8].code+"",g,b,f,u;function _(){return l[7](l[8])}return{key:s,first:null,c(){a=n("button"),g=p(i),b=d(),m(a,"class","tab-item"),re(a,"active",l[3]===l[8].code),this.first=a},m(R,A){r(R,a,A),t(a,g),t(a,b),f||(u=Tt(a,"click",_),f=!0)},p(R,A){l=R,A&16&&i!==(i=l[8].code+"")&&De(g,i),A&24&&re(a,"active",l[3]===l[8].code)},d(R){R&&c(a),f=!1,u()}}}function St(s,l){let a,i,g,b;return i=new vt({props:{content:l[8].body}}),{key:s,first:null,c(){a=n("div"),ne(i.$$.fragment),g=d(),m(a,"class","tab-item"),re(a,"active",l[3]===l[8].code),this.first=a},m(f,u){r(f,a,u),se(i,a,null),t(a,g),b=!0},p(f,u){l=f;const _={};u&16&&(_.content=l[8].body),i.$set(_),(!b||u&24)&&re(a,"active",l[3]===l[8].code)},i(f){b||(Z(i.$$.fragment,f),b=!0)},o(f){x(i.$$.fragment,f),b=!1},d(f){f&&c(a),ie(i)}}}function Lt(s){var rt,ct;let l,a,i=s[0].name+"",g,b,f,u,_,R,A,C,q,Ee,ce,T,de,N,ue,U,ee,We,te,I,Le,pe,le=s[0].name+"",fe,qe,he,V,be,M,me,Be,Q,D,_e,Fe,ke,He,$,Ye,ge,Se,ve,Ne,we,ye,j,$e,E,Pe,Ie,J,W,Re,Ve,Ae,Qe,k,je,B,Je,Ke,ze,Ce,Ge,Oe,Xe,Ze,xe,Te,et,tt,F,Ue,K,Me,L,z,O=[],lt=new Map,at,G,S=[],ot=new Map,H;function nt(e,o){if(e[1]&&e[2])return Wt;if(e[1])return Et;if(e[2])return Dt}let Y=nt(s),P=Y&&Y(s);T=new Ut({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${s[6]}');

        ...

        const authData = await pb.collection('${(rt=s[0])==null?void 0:rt.name}').authWithPassword(
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

        final authData = await pb.collection('${(ct=s[0])==null?void 0:ct.name}').authWithPassword(
          '${s[5]}',
          'YOUR_PASSWORD',
        );

        // after the above you can also access the auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.model.id);

        // "logout" the last authenticated account
        pb.authStore.clear();
    `}});let v=s[1]&&mt(),w=s[1]&&s[2]&&_t(),y=s[2]&&kt();B=new vt({props:{content:"?expand=relField1,relField2.subRelField"}}),F=new Mt({props:{prefix:"record."}});let ae=oe(s[4]);const st=e=>e[8].code;for(let e=0;e<ae.length;e+=1){let o=bt(s,ae,e),h=st(o);lt.set(h,O[e]=gt(h,o))}let X=oe(s[4]);const it=e=>e[8].code;for(let e=0;e<X.length;e+=1){let o=ht(s,X,e),h=it(o);ot.set(h,S[e]=St(h,o))}return{c(){l=n("h3"),a=p("Auth with password ("),g=p(i),b=p(")"),f=d(),u=n("div"),_=n("p"),R=p(`Returns new auth token and account data by a combination of
        `),A=n("strong"),P&&P.c(),C=p(`
        and `),q=n("strong"),q.textContent="password",Ee=p("."),ce=d(),ne(T.$$.fragment),de=d(),N=n("h6"),N.textContent="API details",ue=d(),U=n("div"),ee=n("strong"),ee.textContent="POST",We=d(),te=n("div"),I=n("p"),Le=p("/api/collections/"),pe=n("strong"),fe=p(le),qe=p("/auth-with-password"),he=d(),V=n("div"),V.textContent="Body Parameters",be=d(),M=n("table"),me=n("thead"),me.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr>',Be=d(),Q=n("tbody"),D=n("tr"),_e=n("td"),_e.innerHTML='<div class="inline-flex"><span class="label label-success">Required</span> <span>identity</span></div>',Fe=d(),ke=n("td"),ke.innerHTML='<span class="label">String</span>',He=d(),$=n("td"),Ye=p(`The
                `),v&&v.c(),ge=d(),w&&w.c(),Se=d(),y&&y.c(),ve=p(`
                of the record to authenticate.`),Ne=d(),we=n("tr"),we.innerHTML='<td><div class="inline-flex"><span class="label label-success">Required</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>The auth record password.</td>',ye=d(),j=n("div"),j.textContent="Query parameters",$e=d(),E=n("table"),Pe=n("thead"),Pe.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',Ie=d(),J=n("tbody"),W=n("tr"),Re=n("td"),Re.textContent="expand",Ve=d(),Ae=n("td"),Ae.innerHTML='<span class="label">String</span>',Qe=d(),k=n("td"),je=p(`Auto expand record relations. Ex.:
                `),ne(B.$$.fragment),Je=p(`
                Supports up to 6-levels depth nested relations expansion. `),Ke=n("br"),ze=p(`
                The expanded relations will be appended to the record under the
                `),Ce=n("code"),Ce.textContent="expand",Ge=p(" property (eg. "),Oe=n("code"),Oe.textContent='"expand": {"relField1": {...}, ...}',Xe=p(`).
                `),Ze=n("br"),xe=p(`
                Only the relations to which the request user has permissions to `),Te=n("strong"),Te.textContent="view",et=p(" will be expanded."),tt=d(),ne(F.$$.fragment),Ue=d(),K=n("div"),K.textContent="Responses",Me=d(),L=n("div"),z=n("div");for(let e=0;e<O.length;e+=1)O[e].c();at=d(),G=n("div");for(let e=0;e<S.length;e+=1)S[e].c();m(l,"class","m-b-sm"),m(u,"class","content txt-lg m-b-sm"),m(N,"class","m-b-xs"),m(ee,"class","label label-primary"),m(te,"class","content"),m(U,"class","alert alert-success"),m(V,"class","section-title"),m(M,"class","table-compact table-border m-b-base"),m(j,"class","section-title"),m(E,"class","table-compact table-border m-b-base"),m(K,"class","section-title"),m(z,"class","tabs-header compact combined left"),m(G,"class","tabs-content"),m(L,"class","tabs")},m(e,o){r(e,l,o),t(l,a),t(l,g),t(l,b),r(e,f,o),r(e,u,o),t(u,_),t(_,R),t(_,A),P&&P.m(A,null),t(_,C),t(_,q),t(_,Ee),r(e,ce,o),se(T,e,o),r(e,de,o),r(e,N,o),r(e,ue,o),r(e,U,o),t(U,ee),t(U,We),t(U,te),t(te,I),t(I,Le),t(I,pe),t(pe,fe),t(I,qe),r(e,he,o),r(e,V,o),r(e,be,o),r(e,M,o),t(M,me),t(M,Be),t(M,Q),t(Q,D),t(D,_e),t(D,Fe),t(D,ke),t(D,He),t(D,$),t($,Ye),v&&v.m($,null),t($,ge),w&&w.m($,null),t($,Se),y&&y.m($,null),t($,ve),t(Q,Ne),t(Q,we),r(e,ye,o),r(e,j,o),r(e,$e,o),r(e,E,o),t(E,Pe),t(E,Ie),t(E,J),t(J,W),t(W,Re),t(W,Ve),t(W,Ae),t(W,Qe),t(W,k),t(k,je),se(B,k,null),t(k,Je),t(k,Ke),t(k,ze),t(k,Ce),t(k,Ge),t(k,Oe),t(k,Xe),t(k,Ze),t(k,xe),t(k,Te),t(k,et),t(J,tt),se(F,J,null),r(e,Ue,o),r(e,K,o),r(e,Me,o),r(e,L,o),t(L,z);for(let h=0;h<O.length;h+=1)O[h]&&O[h].m(z,null);t(L,at),t(L,G);for(let h=0;h<S.length;h+=1)S[h]&&S[h].m(G,null);H=!0},p(e,[o]){var dt,ut;(!H||o&1)&&i!==(i=e[0].name+"")&&De(g,i),Y!==(Y=nt(e))&&(P&&P.d(1),P=Y&&Y(e),P&&(P.c(),P.m(A,null)));const h={};o&97&&(h.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[6]}');

        ...

        const authData = await pb.collection('${(dt=e[0])==null?void 0:dt.name}').authWithPassword(
            '${e[5]}',
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

        final pb = PocketBase('${e[6]}');

        ...

        final authData = await pb.collection('${(ut=e[0])==null?void 0:ut.name}').authWithPassword(
          '${e[5]}',
          'YOUR_PASSWORD',
        );

        // after the above you can also access the auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.model.id);

        // "logout" the last authenticated account
        pb.authStore.clear();
    `),T.$set(h),(!H||o&1)&&le!==(le=e[0].name+"")&&De(fe,le),e[1]?v||(v=mt(),v.c(),v.m($,ge)):v&&(v.d(1),v=null),e[1]&&e[2]?w||(w=_t(),w.c(),w.m($,Se)):w&&(w.d(1),w=null),e[2]?y||(y=kt(),y.c(),y.m($,ve)):y&&(y.d(1),y=null),o&24&&(ae=oe(e[4]),O=pt(O,o,st,1,e,ae,lt,z,Pt,gt,null,bt)),o&24&&(X=oe(e[4]),Rt(),S=pt(S,o,it,1,e,X,ot,G,At,St,null,ht),Ct())},i(e){if(!H){Z(T.$$.fragment,e),Z(B.$$.fragment,e),Z(F.$$.fragment,e);for(let o=0;o<X.length;o+=1)Z(S[o]);H=!0}},o(e){x(T.$$.fragment,e),x(B.$$.fragment,e),x(F.$$.fragment,e);for(let o=0;o<S.length;o+=1)x(S[o]);H=!1},d(e){e&&(c(l),c(f),c(u),c(ce),c(de),c(N),c(ue),c(U),c(he),c(V),c(be),c(M),c(ye),c(j),c($e),c(E),c(Ue),c(K),c(Me),c(L)),P&&P.d(),ie(T,e),v&&v.d(),w&&w.d(),y&&y.d(),ie(B),ie(F);for(let o=0;o<O.length;o+=1)O[o].d();for(let o=0;o<S.length;o+=1)S[o].d()}}}function qt(s,l,a){let i,g,b,f,{collection:u}=l,_=200,R=[];const A=C=>a(3,_=C.code);return s.$$set=C=>{"collection"in C&&a(0,u=C.collection)},s.$$.update=()=>{var C,q;s.$$.dirty&1&&a(2,g=(C=u==null?void 0:u.options)==null?void 0:C.allowEmailAuth),s.$$.dirty&1&&a(1,b=(q=u==null?void 0:u.options)==null?void 0:q.allowUsernameAuth),s.$$.dirty&6&&a(5,f=b&&g?"YOUR_USERNAME_OR_EMAIL":b?"YOUR_USERNAME":"YOUR_EMAIL"),s.$$.dirty&1&&a(4,R=[{code:200,body:JSON.stringify({token:"JWT_TOKEN",record:ft.dummyCollectionRecord(u)},null,2)},{code:400,body:`
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
            `}])},a(6,i=ft.getApiExampleUrl(Ot.baseUrl)),[u,b,g,_,R,f,i,A]}class Yt extends wt{constructor(l){super(),yt(this,l,qt,Lt,$t,{collection:0})}}export{Yt as default};
