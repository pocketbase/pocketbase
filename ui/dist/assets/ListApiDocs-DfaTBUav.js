import{S as el,i as ll,s as sl,H as ze,h as m,l as h,o as nl,u as e,v as s,L as ol,w as a,n as t,A as g,V as al,W as Le,X as ae,d as Kt,Y as il,t as Ct,a as kt,I as ve,Z as Je,_ as rl,C as cl,$ as dl,D as pl,m as Qt,c as Vt,J as Te,p as fl,k as Ae}from"./index-CkwOC79g.js";import{F as ul}from"./FieldsQueryParam-haBHHZQ1.js";function ml(r){let n,o,i;return{c(){n=e("span"),n.textContent="Show details",o=s(),i=e("i"),a(n,"class","txt"),a(i,"class","ri-arrow-down-s-line")},m(f,b){h(f,n,b),h(f,o,b),h(f,i,b)},d(f){f&&(m(n),m(o),m(i))}}}function hl(r){let n,o,i;return{c(){n=e("span"),n.textContent="Hide details",o=s(),i=e("i"),a(n,"class","txt"),a(i,"class","ri-arrow-up-s-line")},m(f,b){h(f,n,b),h(f,o,b),h(f,i,b)},d(f){f&&(m(n),m(o),m(i))}}}function Ke(r){let n,o,i,f,b,p,u,C,_,x,d,Y,yt,Wt,E,Xt,D,it,P,Z,ie,j,U,re,rt,vt,tt,Ft,ce,ct,dt,et,N,Yt,Lt,k,lt,At,Zt,Tt,z,st,Pt,te,Rt,v,pt,Ot,de,ft,pe,H,St,nt,Et,F,ut,fe,J,Nt,ee,qt,le,Dt,ue,L,mt,me,ht,he,M,be,T,Ht,ot,Mt,K,bt,ge,I,It,y,Bt,at,Gt,_e,Q,gt,we,_t,xe,jt,$e,B,Ut,Ce,G,ke,wt,se,R,xt,V,W,O,zt,ne,X;return{c(){n=e("p"),n.innerHTML=`The syntax basically follows the format
        <code><span class="txt-success">OPERAND</span> <span class="txt-danger">OPERATOR</span> <span class="txt-success">OPERAND</span></code>, where:`,o=s(),i=e("ul"),f=e("li"),f.innerHTML=`<code class="txt-success">OPERAND</code> - could be any of the above field literal, string (single
            or double quoted), number, null, true, false`,b=s(),p=e("li"),u=e("code"),u.textContent="OPERATOR",C=g(` - is one of:
            `),_=e("br"),x=s(),d=e("ul"),Y=e("li"),yt=e("code"),yt.textContent="=",Wt=s(),E=e("span"),E.textContent="Equal",Xt=s(),D=e("li"),it=e("code"),it.textContent="!=",P=s(),Z=e("span"),Z.textContent="NOT equal",ie=s(),j=e("li"),U=e("code"),U.textContent=">",re=s(),rt=e("span"),rt.textContent="Greater than",vt=s(),tt=e("li"),Ft=e("code"),Ft.textContent=">=",ce=s(),ct=e("span"),ct.textContent="Greater than or equal",dt=s(),et=e("li"),N=e("code"),N.textContent="<",Yt=s(),Lt=e("span"),Lt.textContent="Less than",k=s(),lt=e("li"),At=e("code"),At.textContent="<=",Zt=s(),Tt=e("span"),Tt.textContent="Less than or equal",z=s(),st=e("li"),Pt=e("code"),Pt.textContent="~",te=s(),Rt=e("span"),Rt.textContent=`Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,v=s(),pt=e("li"),Ot=e("code"),Ot.textContent="!~",de=s(),ft=e("span"),ft.textContent=`NOT Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,pe=s(),H=e("li"),St=e("code"),St.textContent="?=",nt=s(),Et=e("em"),Et.textContent="Any/At least one of",F=s(),ut=e("span"),ut.textContent="Equal",fe=s(),J=e("li"),Nt=e("code"),Nt.textContent="?!=",ee=s(),qt=e("em"),qt.textContent="Any/At least one of",le=s(),Dt=e("span"),Dt.textContent="NOT equal",ue=s(),L=e("li"),mt=e("code"),mt.textContent="?>",me=s(),ht=e("em"),ht.textContent="Any/At least one of",he=s(),M=e("span"),M.textContent="Greater than",be=s(),T=e("li"),Ht=e("code"),Ht.textContent="?>=",ot=s(),Mt=e("em"),Mt.textContent="Any/At least one of",K=s(),bt=e("span"),bt.textContent="Greater than or equal",ge=s(),I=e("li"),It=e("code"),It.textContent="?<",y=s(),Bt=e("em"),Bt.textContent="Any/At least one of",at=s(),Gt=e("span"),Gt.textContent="Less than",_e=s(),Q=e("li"),gt=e("code"),gt.textContent="?<=",we=s(),_t=e("em"),_t.textContent="Any/At least one of",xe=s(),jt=e("span"),jt.textContent="Less than or equal",$e=s(),B=e("li"),Ut=e("code"),Ut.textContent="?~",Ce=s(),G=e("em"),G.textContent="Any/At least one of",ke=s(),wt=e("span"),wt.textContent=`Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,se=s(),R=e("li"),xt=e("code"),xt.textContent="?!~",V=s(),W=e("em"),W.textContent="Any/At least one of",O=s(),zt=e("span"),zt.textContent=`NOT Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,ne=s(),X=e("p"),X.innerHTML=`To group and combine several expressions you could use brackets
        <code>(...)</code>, <code>&amp;&amp;</code> (AND) and <code>||</code> (OR) tokens.`,a(u,"class","txt-danger"),a(yt,"class","filter-op svelte-1w7s5nw"),a(E,"class","txt"),a(it,"class","filter-op svelte-1w7s5nw"),a(Z,"class","txt"),a(U,"class","filter-op svelte-1w7s5nw"),a(rt,"class","txt"),a(Ft,"class","filter-op svelte-1w7s5nw"),a(ct,"class","txt"),a(N,"class","filter-op svelte-1w7s5nw"),a(Lt,"class","txt"),a(At,"class","filter-op svelte-1w7s5nw"),a(Tt,"class","txt"),a(Pt,"class","filter-op svelte-1w7s5nw"),a(Rt,"class","txt"),a(Ot,"class","filter-op svelte-1w7s5nw"),a(ft,"class","txt"),a(St,"class","filter-op svelte-1w7s5nw"),a(Et,"class","txt-hint"),a(ut,"class","txt"),a(Nt,"class","filter-op svelte-1w7s5nw"),a(qt,"class","txt-hint"),a(Dt,"class","txt"),a(mt,"class","filter-op svelte-1w7s5nw"),a(ht,"class","txt-hint"),a(M,"class","txt"),a(Ht,"class","filter-op svelte-1w7s5nw"),a(Mt,"class","txt-hint"),a(bt,"class","txt"),a(It,"class","filter-op svelte-1w7s5nw"),a(Bt,"class","txt-hint"),a(Gt,"class","txt"),a(gt,"class","filter-op svelte-1w7s5nw"),a(_t,"class","txt-hint"),a(jt,"class","txt"),a(Ut,"class","filter-op svelte-1w7s5nw"),a(G,"class","txt-hint"),a(wt,"class","txt"),a(xt,"class","filter-op svelte-1w7s5nw"),a(W,"class","txt-hint"),a(zt,"class","txt")},m($,$t){h($,n,$t),h($,o,$t),h($,i,$t),t(i,f),t(i,b),t(i,p),t(p,u),t(p,C),t(p,_),t(p,x),t(p,d),t(d,Y),t(Y,yt),t(Y,Wt),t(Y,E),t(d,Xt),t(d,D),t(D,it),t(D,P),t(D,Z),t(d,ie),t(d,j),t(j,U),t(j,re),t(j,rt),t(d,vt),t(d,tt),t(tt,Ft),t(tt,ce),t(tt,ct),t(d,dt),t(d,et),t(et,N),t(et,Yt),t(et,Lt),t(d,k),t(d,lt),t(lt,At),t(lt,Zt),t(lt,Tt),t(d,z),t(d,st),t(st,Pt),t(st,te),t(st,Rt),t(d,v),t(d,pt),t(pt,Ot),t(pt,de),t(pt,ft),t(d,pe),t(d,H),t(H,St),t(H,nt),t(H,Et),t(H,F),t(H,ut),t(d,fe),t(d,J),t(J,Nt),t(J,ee),t(J,qt),t(J,le),t(J,Dt),t(d,ue),t(d,L),t(L,mt),t(L,me),t(L,ht),t(L,he),t(L,M),t(d,be),t(d,T),t(T,Ht),t(T,ot),t(T,Mt),t(T,K),t(T,bt),t(d,ge),t(d,I),t(I,It),t(I,y),t(I,Bt),t(I,at),t(I,Gt),t(d,_e),t(d,Q),t(Q,gt),t(Q,we),t(Q,_t),t(Q,xe),t(Q,jt),t(d,$e),t(d,B),t(B,Ut),t(B,Ce),t(B,G),t(B,ke),t(B,wt),t(d,se),t(d,R),t(R,xt),t(R,V),t(R,W),t(R,O),t(R,zt),h($,ne,$t),h($,X,$t)},d($){$&&(m(n),m(o),m(i),m(ne),m(X))}}}function bl(r){let n,o,i,f,b;function p(x,d){return x[0]?hl:ml}let u=p(r),C=u(r),_=r[0]&&Ke();return{c(){n=e("button"),C.c(),o=s(),_&&_.c(),i=ol(),a(n,"class","btn btn-sm btn-secondary m-t-10")},m(x,d){h(x,n,d),C.m(n,null),h(x,o,d),_&&_.m(x,d),h(x,i,d),f||(b=nl(n,"click",r[1]),f=!0)},p(x,[d]){u!==(u=p(x))&&(C.d(1),C=u(x),C&&(C.c(),C.m(n,null))),x[0]?_||(_=Ke(),_.c(),_.m(i.parentNode,i)):_&&(_.d(1),_=null)},i:ze,o:ze,d(x){x&&(m(n),m(o),m(i)),C.d(),_&&_.d(x),f=!1,b()}}}function gl(r,n,o){let i=!1;function f(){o(0,i=!i)}return[i,f]}class _l extends el{constructor(n){super(),ll(this,n,gl,bl,sl,{})}}function Qe(r,n,o){const i=r.slice();return i[8]=n[o],i}function Ve(r,n,o){const i=r.slice();return i[8]=n[o],i}function We(r,n,o){const i=r.slice();return i[13]=n[o],i[15]=o,i}function Xe(r){let n;return{c(){n=e("p"),n.innerHTML="Requires superuser <code>Authorization:TOKEN</code> header",a(n,"class","txt-hint txt-sm txt-right")},m(o,i){h(o,n,i)},d(o){o&&m(n)}}}function Ye(r){let n,o=r[13]+"",i,f=r[15]<r[4].length-1?", ":"",b;return{c(){n=e("code"),i=g(o),b=g(f)},m(p,u){h(p,n,u),t(n,i),h(p,b,u)},p(p,u){u&16&&o!==(o=p[13]+"")&&ve(i,o),u&16&&f!==(f=p[15]<p[4].length-1?", ":"")&&ve(b,f)},d(p){p&&(m(n),m(b))}}}function Ze(r,n){let o,i,f;function b(){return n[7](n[8])}return{key:r,first:null,c(){o=e("button"),o.textContent=`${n[8].code} `,a(o,"type","button"),a(o,"class","tab-item"),Ae(o,"active",n[2]===n[8].code),this.first=o},m(p,u){h(p,o,u),i||(f=nl(o,"click",b),i=!0)},p(p,u){n=p,u&36&&Ae(o,"active",n[2]===n[8].code)},d(p){p&&m(o),i=!1,f()}}}function tl(r,n){let o,i,f,b;return i=new Le({props:{content:n[8].body}}),{key:r,first:null,c(){o=e("div"),Vt(i.$$.fragment),f=s(),a(o,"class","tab-item"),Ae(o,"active",n[2]===n[8].code),this.first=o},m(p,u){h(p,o,u),Qt(i,o,null),t(o,f),b=!0},p(p,u){n=p,(!b||u&36)&&Ae(o,"active",n[2]===n[8].code)},i(p){b||(kt(i.$$.fragment,p),b=!0)},o(p){Ct(i.$$.fragment,p),b=!1},d(p){p&&m(o),Kt(i)}}}function wl(r){var Oe,Se,Ee,Ne,qe,De;let n,o,i=r[0].name+"",f,b,p,u,C,_,x,d=r[0].name+"",Y,yt,Wt,E,Xt,D,it,P,Z,ie,j,U,re,rt,vt=r[0].name+"",tt,Ft,ce,ct,dt,et,N,Yt,Lt,k,lt,At,Zt,Tt,z,st,Pt,te,Rt,v,pt,Ot,de,ft,pe,H,St,nt,Et,F,ut,fe,J,Nt,ee,qt,le,Dt,ue,L,mt,me,ht,he,M,be,T,Ht,ot,Mt,K,bt,ge,I,It,y,Bt,at,Gt,_e,Q,gt,we,_t,xe,jt,$e,B,Ut,Ce,G,ke,wt,se,R,xt,V,W,O=[],zt=new Map,ne,X,$=[],$t=new Map,Jt;E=new al({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${r[3]}');

        ...

        // fetch a paginated records list
        const resultList = await pb.collection('${(Oe=r[0])==null?void 0:Oe.name}').getList(1, 50, {
            filter: 'someField1 != someField2',
        });

        // you can also fetch all records at once via getFullList
        const records = await pb.collection('${(Se=r[0])==null?void 0:Se.name}').getFullList({
            sort: '-someField',
        });

        // or fetch only the first record that matches the specified filter
        const record = await pb.collection('${(Ee=r[0])==null?void 0:Ee.name}').getFirstListItem('someField="test"', {
            expand: 'relField1,relField2.subRelField',
        });
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${r[3]}');

        ...

        // fetch a paginated records list
        final resultList = await pb.collection('${(Ne=r[0])==null?void 0:Ne.name}').getList(
          page: 1,
          perPage: 50,
          filter: 'someField1 != someField2',
        );

        // you can also fetch all records at once via getFullList
        final records = await pb.collection('${(qe=r[0])==null?void 0:qe.name}').getFullList(
          sort: '-someField',
        );

        // or fetch only the first record that matches the specified filter
        final record = await pb.collection('${(De=r[0])==null?void 0:De.name}').getFirstListItem(
          'someField="test"',
          expand: 'relField1,relField2.subRelField',
        );
    `}});let S=r[1]&&Xe();nt=new Le({props:{content:`
                        // DESC by created and ASC by id
                        ?sort=-created,id
                    `}});let oe=ae(r[4]),A=[];for(let l=0;l<oe.length;l+=1)A[l]=Ye(We(r,oe,l));T=new Le({props:{content:`
                        ?filter=(id='abc' && created>'2022-01-01')
                    `}}),ot=new _l({}),at=new Le({props:{content:"?expand=relField1,relField2.subRelField"}}),G=new ul({});let Fe=ae(r[5]);const Pe=l=>l[8].code;for(let l=0;l<Fe.length;l+=1){let c=Ve(r,Fe,l),w=Pe(c);zt.set(w,O[l]=Ze(w,c))}let ye=ae(r[5]);const Re=l=>l[8].code;for(let l=0;l<ye.length;l+=1){let c=Qe(r,ye,l),w=Re(c);$t.set(w,$[l]=tl(w,c))}return{c(){n=e("h3"),o=g("List/Search ("),f=g(i),b=g(")"),p=s(),u=e("div"),C=e("p"),_=g("Fetch a paginated "),x=e("strong"),Y=g(d),yt=g(" records list, supporting sorting and filtering."),Wt=s(),Vt(E.$$.fragment),Xt=s(),D=e("h6"),D.textContent="API details",it=s(),P=e("div"),Z=e("strong"),Z.textContent="GET",ie=s(),j=e("div"),U=e("p"),re=g("/api/collections/"),rt=e("strong"),tt=g(vt),Ft=g("/records"),ce=s(),S&&S.c(),ct=s(),dt=e("div"),dt.textContent="Query parameters",et=s(),N=e("table"),Yt=e("thead"),Yt.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',Lt=s(),k=e("tbody"),lt=e("tr"),lt.innerHTML='<td>page</td> <td><span class="label">Number</span></td> <td>The page (aka. offset) of the paginated list (default to 1).</td>',At=s(),Zt=e("tr"),Zt.innerHTML='<td>perPage</td> <td><span class="label">Number</span></td> <td>Specify the max returned records per page (default to 30).</td>',Tt=s(),z=e("tr"),st=e("td"),st.textContent="sort",Pt=s(),te=e("td"),te.innerHTML='<span class="label">String</span>',Rt=s(),v=e("td"),pt=g("Specify the records order attribute(s). "),Ot=e("br"),de=g(`
                Add `),ft=e("code"),ft.textContent="-",pe=g(" / "),H=e("code"),H.textContent="+",St=g(` (default) in front of the attribute for DESC / ASC order.
                Ex.:
                `),Vt(nt.$$.fragment),Et=s(),F=e("p"),ut=e("strong"),ut.textContent="Supported record sort fields:",fe=s(),J=e("br"),Nt=s(),ee=e("code"),ee.textContent="@random",qt=g(`,
                    `),le=e("code"),le.textContent="@rowid",Dt=g(`,
                    `);for(let l=0;l<A.length;l+=1)A[l].c();ue=s(),L=e("tr"),mt=e("td"),mt.textContent="filter",me=s(),ht=e("td"),ht.innerHTML='<span class="label">String</span>',he=s(),M=e("td"),be=g(`Filter the returned records. Ex.:
                `),Vt(T.$$.fragment),Ht=s(),Vt(ot.$$.fragment),Mt=s(),K=e("tr"),bt=e("td"),bt.textContent="expand",ge=s(),I=e("td"),I.innerHTML='<span class="label">String</span>',It=s(),y=e("td"),Bt=g(`Auto expand record relations. Ex.:
                `),Vt(at.$$.fragment),Gt=g(`
                Supports up to 6-levels depth nested relations expansion. `),_e=e("br"),Q=g(`
                The expanded relations will be appended to each individual record under the
                `),gt=e("code"),gt.textContent="expand",we=g(" property (eg. "),_t=e("code"),_t.textContent='"expand": {"relField1": {...}, ...}',xe=g(`).
                `),jt=e("br"),$e=g(`
                Only the relations to which the request user has permissions to `),B=e("strong"),B.textContent="view",Ut=g(" will be expanded."),Ce=s(),Vt(G.$$.fragment),ke=s(),wt=e("tr"),wt.innerHTML=`<td id="query-page">skipTotal</td> <td><span class="label">Boolean</span></td> <td>If it is set the total counts query will be skipped and the response fields
                <code>totalItems</code> and <code>totalPages</code> will have <code>-1</code> value.
                <br/>
                This could drastically speed up the search queries when the total counters are not needed or cursor
                based pagination is used.
                <br/>
                For optimization purposes, it is set by default for the
                <code>getFirstListItem()</code>
                and
                <code>getFullList()</code> SDKs methods.</td>`,se=s(),R=e("div"),R.textContent="Responses",xt=s(),V=e("div"),W=e("div");for(let l=0;l<O.length;l+=1)O[l].c();ne=s(),X=e("div");for(let l=0;l<$.length;l+=1)$[l].c();a(n,"class","m-b-sm"),a(u,"class","content txt-lg m-b-sm"),a(D,"class","m-b-xs"),a(Z,"class","label label-primary"),a(j,"class","content"),a(P,"class","alert alert-info"),a(dt,"class","section-title"),a(N,"class","table-compact table-border m-b-base"),a(R,"class","section-title"),a(W,"class","tabs-header compact combined left"),a(X,"class","tabs-content"),a(V,"class","tabs")},m(l,c){h(l,n,c),t(n,o),t(n,f),t(n,b),h(l,p,c),h(l,u,c),t(u,C),t(C,_),t(C,x),t(x,Y),t(C,yt),h(l,Wt,c),Qt(E,l,c),h(l,Xt,c),h(l,D,c),h(l,it,c),h(l,P,c),t(P,Z),t(P,ie),t(P,j),t(j,U),t(U,re),t(U,rt),t(rt,tt),t(U,Ft),t(P,ce),S&&S.m(P,null),h(l,ct,c),h(l,dt,c),h(l,et,c),h(l,N,c),t(N,Yt),t(N,Lt),t(N,k),t(k,lt),t(k,At),t(k,Zt),t(k,Tt),t(k,z),t(z,st),t(z,Pt),t(z,te),t(z,Rt),t(z,v),t(v,pt),t(v,Ot),t(v,de),t(v,ft),t(v,pe),t(v,H),t(v,St),Qt(nt,v,null),t(v,Et),t(v,F),t(F,ut),t(F,fe),t(F,J),t(F,Nt),t(F,ee),t(F,qt),t(F,le),t(F,Dt);for(let w=0;w<A.length;w+=1)A[w]&&A[w].m(F,null);t(k,ue),t(k,L),t(L,mt),t(L,me),t(L,ht),t(L,he),t(L,M),t(M,be),Qt(T,M,null),t(M,Ht),Qt(ot,M,null),t(k,Mt),t(k,K),t(K,bt),t(K,ge),t(K,I),t(K,It),t(K,y),t(y,Bt),Qt(at,y,null),t(y,Gt),t(y,_e),t(y,Q),t(y,gt),t(y,we),t(y,_t),t(y,xe),t(y,jt),t(y,$e),t(y,B),t(y,Ut),t(k,Ce),Qt(G,k,null),t(k,ke),t(k,wt),h(l,se,c),h(l,R,c),h(l,xt,c),h(l,V,c),t(V,W);for(let w=0;w<O.length;w+=1)O[w]&&O[w].m(W,null);t(V,ne),t(V,X);for(let w=0;w<$.length;w+=1)$[w]&&$[w].m(X,null);Jt=!0},p(l,[c]){var He,Me,Ie,Be,Ge,je;(!Jt||c&1)&&i!==(i=l[0].name+"")&&ve(f,i),(!Jt||c&1)&&d!==(d=l[0].name+"")&&ve(Y,d);const w={};if(c&9&&(w.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${l[3]}');

        ...

        // fetch a paginated records list
        const resultList = await pb.collection('${(He=l[0])==null?void 0:He.name}').getList(1, 50, {
            filter: 'someField1 != someField2',
        });

        // you can also fetch all records at once via getFullList
        const records = await pb.collection('${(Me=l[0])==null?void 0:Me.name}').getFullList({
            sort: '-someField',
        });

        // or fetch only the first record that matches the specified filter
        const record = await pb.collection('${(Ie=l[0])==null?void 0:Ie.name}').getFirstListItem('someField="test"', {
            expand: 'relField1,relField2.subRelField',
        });
    `),c&9&&(w.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${l[3]}');

        ...

        // fetch a paginated records list
        final resultList = await pb.collection('${(Be=l[0])==null?void 0:Be.name}').getList(
          page: 1,
          perPage: 50,
          filter: 'someField1 != someField2',
        );

        // you can also fetch all records at once via getFullList
        final records = await pb.collection('${(Ge=l[0])==null?void 0:Ge.name}').getFullList(
          sort: '-someField',
        );

        // or fetch only the first record that matches the specified filter
        final record = await pb.collection('${(je=l[0])==null?void 0:je.name}').getFirstListItem(
          'someField="test"',
          expand: 'relField1,relField2.subRelField',
        );
    `),E.$set(w),(!Jt||c&1)&&vt!==(vt=l[0].name+"")&&ve(tt,vt),l[1]?S||(S=Xe(),S.c(),S.m(P,null)):S&&(S.d(1),S=null),c&16){oe=ae(l[4]);let q;for(q=0;q<oe.length;q+=1){const Ue=We(l,oe,q);A[q]?A[q].p(Ue,c):(A[q]=Ye(Ue),A[q].c(),A[q].m(F,null))}for(;q<A.length;q+=1)A[q].d(1);A.length=oe.length}c&36&&(Fe=ae(l[5]),O=Je(O,c,Pe,1,l,Fe,zt,W,rl,Ze,null,Ve)),c&36&&(ye=ae(l[5]),cl(),$=Je($,c,Re,1,l,ye,$t,X,dl,tl,null,Qe),pl())},i(l){if(!Jt){kt(E.$$.fragment,l),kt(nt.$$.fragment,l),kt(T.$$.fragment,l),kt(ot.$$.fragment,l),kt(at.$$.fragment,l),kt(G.$$.fragment,l);for(let c=0;c<ye.length;c+=1)kt($[c]);Jt=!0}},o(l){Ct(E.$$.fragment,l),Ct(nt.$$.fragment,l),Ct(T.$$.fragment,l),Ct(ot.$$.fragment,l),Ct(at.$$.fragment,l),Ct(G.$$.fragment,l);for(let c=0;c<$.length;c+=1)Ct($[c]);Jt=!1},d(l){l&&(m(n),m(p),m(u),m(Wt),m(Xt),m(D),m(it),m(P),m(ct),m(dt),m(et),m(N),m(se),m(R),m(xt),m(V)),Kt(E,l),S&&S.d(),Kt(nt),il(A,l),Kt(T),Kt(ot),Kt(at),Kt(G);for(let c=0;c<O.length;c+=1)O[c].d();for(let c=0;c<$.length;c+=1)$[c].d()}}}function xl(r,n,o){let i,f,b,p,{collection:u}=n,C=200,_=[];const x=d=>o(2,C=d.code);return r.$$set=d=>{"collection"in d&&o(0,u=d.collection)},r.$$.update=()=>{r.$$.dirty&1&&o(4,i=Te.getAllCollectionIdentifiers(u)),r.$$.dirty&1&&o(1,f=(u==null?void 0:u.listRule)===null),r.$$.dirty&1&&o(6,p=Te.dummyCollectionRecord(u)),r.$$.dirty&67&&u!=null&&u.id&&(_.push({code:200,body:JSON.stringify({page:1,perPage:30,totalPages:1,totalItems:2,items:[p,Object.assign({},p,{id:p.id+"2"})]},null,2)}),_.push({code:400,body:`
                {
                  "status": 400,
                  "message": "Something went wrong while processing your request. Invalid filter.",
                  "data": {}
                }
            `}),f&&_.push({code:403,body:`
                    {
                      "status": 403,
                      "message": "Only superusers can access this action.",
                      "data": {}
                    }
                `}))},o(3,b=Te.getApiExampleUrl(fl.baseURL)),[u,f,C,b,i,_,p,x]}class kl extends el{constructor(n){super(),ll(this,n,xl,wl,sl,{collection:0})}}export{kl as default};
