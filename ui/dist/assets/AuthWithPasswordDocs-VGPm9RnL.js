import{S as kt,i as gt,s as vt,V as St,X as L,W as _t,h as c,d as ae,Y as wt,t as X,a as Z,I as z,Z as ct,_ as yt,C as $t,$ as Pt,D as Ct,l as d,n as t,m as oe,u as s,A as f,v as u,c as se,w as k,J as dt,p as Rt,k as ne,o as Ot}from"./index-pGELYd11.js";import{F as Tt}from"./FieldsQueryParam-cCCbyKA6.js";function pt(i,o,a){const n=i.slice();return n[7]=o[a],n}function ut(i,o,a){const n=i.slice();return n[7]=o[a],n}function ht(i,o,a){const n=i.slice();return n[12]=o[a],n[14]=a,n}function At(i){let o;return{c(){o=f("or")},m(a,n){d(a,o,n)},d(a){a&&c(o)}}}function bt(i){let o,a,n=i[12]+"",m,b=i[14]>0&&At();return{c(){b&&b.c(),o=u(),a=s("strong"),m=f(n)},m(r,h){b&&b.m(r,h),d(r,o,h),d(r,a,h),t(a,m)},p(r,h){h&2&&n!==(n=r[12]+"")&&z(m,n)},d(r){r&&(c(o),c(a)),b&&b.d(r)}}}function ft(i,o){let a,n=o[7].code+"",m,b,r,h;function g(){return o[6](o[7])}return{key:i,first:null,c(){a=s("button"),m=f(n),b=u(),k(a,"class","tab-item"),ne(a,"active",o[2]===o[7].code),this.first=a},m($,_){d($,a,_),t(a,m),t(a,b),r||(h=Ot(a,"click",g),r=!0)},p($,_){o=$,_&8&&n!==(n=o[7].code+"")&&z(m,n),_&12&&ne(a,"active",o[2]===o[7].code)},d($){$&&c(a),r=!1,h()}}}function mt(i,o){let a,n,m,b;return n=new _t({props:{content:o[7].body}}),{key:i,first:null,c(){a=s("div"),se(n.$$.fragment),m=u(),k(a,"class","tab-item"),ne(a,"active",o[2]===o[7].code),this.first=a},m(r,h){d(r,a,h),oe(n,a,null),t(a,m),b=!0},p(r,h){o=r;const g={};h&8&&(g.content=o[7].body),n.$set(g),(!b||h&12)&&ne(a,"active",o[2]===o[7].code)},i(r){b||(Z(n.$$.fragment,r),b=!0)},o(r){X(n.$$.fragment,r),b=!1},d(r){r&&c(a),ae(n)}}}function Dt(i){var ot,st;let o,a,n=i[0].name+"",m,b,r,h,g,$,_,G=i[1].join("/")+"",ie,De,re,We,ce,C,de,q,pe,R,x,Fe,ee,H,Me,ue,te=i[0].name+"",he,Ue,be,Y,fe,O,me,Be,j,T,_e,Le,ke,qe,V,ge,He,ve,Se,E,we,A,ye,Ye,N,D,$e,je,Pe,Ve,v,Ee,M,Ne,Ie,Je,Ce,Qe,Re,Ke,Xe,Ze,Oe,ze,Ge,U,Te,I,Ae,W,J,P=[],xe=new Map,et,Q,w=[],tt=new Map,F;C=new St({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${i[5]}');

        ...

        const authData = await pb.collection('${(ot=i[0])==null?void 0:ot.name}').authWithPassword(
            '${i[4]}',
            'YOUR_PASSWORD',
        );

        // after the above you can also access the auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.record.id);

        // "logout"
        pb.authStore.clear();
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${i[5]}');

        ...

        final authData = await pb.collection('${(st=i[0])==null?void 0:st.name}').authWithPassword(
          '${i[4]}',
          'YOUR_PASSWORD',
        );

        // after the above you can also access the auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.record.id);

        // "logout"
        pb.authStore.clear();
    `}});let B=L(i[1]),S=[];for(let e=0;e<B.length;e+=1)S[e]=bt(ht(i,B,e));M=new _t({props:{content:"?expand=relField1,relField2.subRelField"}}),U=new Tt({props:{prefix:"record."}});let le=L(i[3]);const lt=e=>e[7].code;for(let e=0;e<le.length;e+=1){let l=ut(i,le,e),p=lt(l);xe.set(p,P[e]=ft(p,l))}let K=L(i[3]);const at=e=>e[7].code;for(let e=0;e<K.length;e+=1){let l=pt(i,K,e),p=at(l);tt.set(p,w[e]=mt(p,l))}return{c(){o=s("h3"),a=f("Auth with password ("),m=f(n),b=f(")"),r=u(),h=s("div"),g=s("p"),$=f(`Authenticate with combination of
        `),_=s("strong"),ie=f(G),De=f(" and "),re=s("strong"),re.textContent="password",We=f("."),ce=u(),se(C.$$.fragment),de=u(),q=s("h6"),q.textContent="API details",pe=u(),R=s("div"),x=s("strong"),x.textContent="POST",Fe=u(),ee=s("div"),H=s("p"),Me=f("/api/collections/"),ue=s("strong"),he=f(te),Ue=f("/auth-with-password"),be=u(),Y=s("div"),Y.textContent="Body Parameters",fe=u(),O=s("table"),me=s("thead"),me.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr>',Be=u(),j=s("tbody"),T=s("tr"),_e=s("td"),_e.innerHTML='<div class="inline-flex"><span class="label label-success">Required</span> <span>identity</span></div>',Le=u(),ke=s("td"),ke.innerHTML='<span class="label">String</span>',qe=u(),V=s("td");for(let e=0;e<S.length;e+=1)S[e].c();ge=f(`
                of the record to authenticate.`),He=u(),ve=s("tr"),ve.innerHTML='<td><div class="inline-flex"><span class="label label-success">Required</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>The auth record password.</td>',Se=u(),E=s("div"),E.textContent="Query parameters",we=u(),A=s("table"),ye=s("thead"),ye.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',Ye=u(),N=s("tbody"),D=s("tr"),$e=s("td"),$e.textContent="expand",je=u(),Pe=s("td"),Pe.innerHTML='<span class="label">String</span>',Ve=u(),v=s("td"),Ee=f(`Auto expand record relations. Ex.:
                `),se(M.$$.fragment),Ne=f(`
                Supports up to 6-levels depth nested relations expansion. `),Ie=s("br"),Je=f(`
                The expanded relations will be appended to the record under the
                `),Ce=s("code"),Ce.textContent="expand",Qe=f(" property (eg. "),Re=s("code"),Re.textContent='"expand": {"relField1": {...}, ...}',Ke=f(`).
                `),Xe=s("br"),Ze=f(`
                Only the relations to which the request user has permissions to `),Oe=s("strong"),Oe.textContent="view",ze=f(" will be expanded."),Ge=u(),se(U.$$.fragment),Te=u(),I=s("div"),I.textContent="Responses",Ae=u(),W=s("div"),J=s("div");for(let e=0;e<P.length;e+=1)P[e].c();et=u(),Q=s("div");for(let e=0;e<w.length;e+=1)w[e].c();k(o,"class","m-b-sm"),k(h,"class","content txt-lg m-b-sm"),k(q,"class","m-b-xs"),k(x,"class","label label-primary"),k(ee,"class","content"),k(R,"class","alert alert-success"),k(Y,"class","section-title"),k(O,"class","table-compact table-border m-b-base"),k(E,"class","section-title"),k(A,"class","table-compact table-border m-b-base"),k(I,"class","section-title"),k(J,"class","tabs-header compact combined left"),k(Q,"class","tabs-content"),k(W,"class","tabs")},m(e,l){d(e,o,l),t(o,a),t(o,m),t(o,b),d(e,r,l),d(e,h,l),t(h,g),t(g,$),t(g,_),t(_,ie),t(g,De),t(g,re),t(g,We),d(e,ce,l),oe(C,e,l),d(e,de,l),d(e,q,l),d(e,pe,l),d(e,R,l),t(R,x),t(R,Fe),t(R,ee),t(ee,H),t(H,Me),t(H,ue),t(ue,he),t(H,Ue),d(e,be,l),d(e,Y,l),d(e,fe,l),d(e,O,l),t(O,me),t(O,Be),t(O,j),t(j,T),t(T,_e),t(T,Le),t(T,ke),t(T,qe),t(T,V);for(let p=0;p<S.length;p+=1)S[p]&&S[p].m(V,null);t(V,ge),t(j,He),t(j,ve),d(e,Se,l),d(e,E,l),d(e,we,l),d(e,A,l),t(A,ye),t(A,Ye),t(A,N),t(N,D),t(D,$e),t(D,je),t(D,Pe),t(D,Ve),t(D,v),t(v,Ee),oe(M,v,null),t(v,Ne),t(v,Ie),t(v,Je),t(v,Ce),t(v,Qe),t(v,Re),t(v,Ke),t(v,Xe),t(v,Ze),t(v,Oe),t(v,ze),t(N,Ge),oe(U,N,null),d(e,Te,l),d(e,I,l),d(e,Ae,l),d(e,W,l),t(W,J);for(let p=0;p<P.length;p+=1)P[p]&&P[p].m(J,null);t(W,et),t(W,Q);for(let p=0;p<w.length;p+=1)w[p]&&w[p].m(Q,null);F=!0},p(e,[l]){var nt,it;(!F||l&1)&&n!==(n=e[0].name+"")&&z(m,n),(!F||l&2)&&G!==(G=e[1].join("/")+"")&&z(ie,G);const p={};if(l&49&&(p.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[5]}');

        ...

        const authData = await pb.collection('${(nt=e[0])==null?void 0:nt.name}').authWithPassword(
            '${e[4]}',
            'YOUR_PASSWORD',
        );

        // after the above you can also access the auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.record.id);

        // "logout"
        pb.authStore.clear();
    `),l&49&&(p.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[5]}');

        ...

        final authData = await pb.collection('${(it=e[0])==null?void 0:it.name}').authWithPassword(
          '${e[4]}',
          'YOUR_PASSWORD',
        );

        // after the above you can also access the auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.record.id);

        // "logout"
        pb.authStore.clear();
    `),C.$set(p),(!F||l&1)&&te!==(te=e[0].name+"")&&z(he,te),l&2){B=L(e[1]);let y;for(y=0;y<B.length;y+=1){const rt=ht(e,B,y);S[y]?S[y].p(rt,l):(S[y]=bt(rt),S[y].c(),S[y].m(V,ge))}for(;y<S.length;y+=1)S[y].d(1);S.length=B.length}l&12&&(le=L(e[3]),P=ct(P,l,lt,1,e,le,xe,J,yt,ft,null,ut)),l&12&&(K=L(e[3]),$t(),w=ct(w,l,at,1,e,K,tt,Q,Pt,mt,null,pt),Ct())},i(e){if(!F){Z(C.$$.fragment,e),Z(M.$$.fragment,e),Z(U.$$.fragment,e);for(let l=0;l<K.length;l+=1)Z(w[l]);F=!0}},o(e){X(C.$$.fragment,e),X(M.$$.fragment,e),X(U.$$.fragment,e);for(let l=0;l<w.length;l+=1)X(w[l]);F=!1},d(e){e&&(c(o),c(r),c(h),c(ce),c(de),c(q),c(pe),c(R),c(be),c(Y),c(fe),c(O),c(Se),c(E),c(we),c(A),c(Te),c(I),c(Ae),c(W)),ae(C,e),wt(S,e),ae(M),ae(U);for(let l=0;l<P.length;l+=1)P[l].d();for(let l=0;l<w.length;l+=1)w[l].d()}}}function Wt(i,o,a){let n,m,b,{collection:r}=o,h=200,g=[];const $=_=>a(2,h=_.code);return i.$$set=_=>{"collection"in _&&a(0,r=_.collection)},i.$$.update=()=>{var _;i.$$.dirty&1&&a(1,m=((_=r==null?void 0:r.passwordAuth)==null?void 0:_.identityFields)||[]),i.$$.dirty&2&&a(4,b=m.length==0?"NONE":"YOUR_"+m.join("_OR_").toUpperCase()),i.$$.dirty&1&&a(3,g=[{code:200,body:JSON.stringify({token:"JWT_TOKEN",record:dt.dummyCollectionRecord(r)},null,2)},{code:400,body:`
                {
                  "status": 400,
                  "message": "Failed to authenticate.",
                  "data": {
                    "identity": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `}])},a(5,n=dt.getApiExampleUrl(Rt.baseURL)),[r,m,h,g,b,n,$]}class Ut extends kt{constructor(o){super(),gt(this,o,Wt,Dt,vt,{collection:0})}}export{Ut as default};
