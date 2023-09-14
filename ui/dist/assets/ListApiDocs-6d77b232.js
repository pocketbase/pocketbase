import{S as Ze,i as tl,s as el,e,b as s,E as sl,f as a,g as u,u as ll,y as Ue,o as m,w as _,h as t,N as Fe,O as te,c as ee,m as le,x as ke,P as je,Q as nl,k as ol,R as al,n as il,t as Bt,a as Gt,d as se,T as rl,C as ve,p as cl,r as Le}from"./index-8354bde7.js";import{S as dl}from"./SdkTabs-86785e52.js";function pl(d){let n,o,i;return{c(){n=e("span"),n.textContent="Show details",o=s(),i=e("i"),a(n,"class","txt"),a(i,"class","ri-arrow-down-s-line")},m(f,h){u(f,n,h),u(f,o,h),u(f,i,h)},d(f){f&&(m(n),m(o),m(i))}}}function fl(d){let n,o,i;return{c(){n=e("span"),n.textContent="Hide details",o=s(),i=e("i"),a(n,"class","txt"),a(i,"class","ri-arrow-up-s-line")},m(f,h){u(f,n,h),u(f,o,h),u(f,i,h)},d(f){f&&(m(n),m(o),m(i))}}}function ze(d){let n,o,i,f,h,r,b,C,$,g,p,Z,Ct,Ut,E,jt,M,it,S,tt,ne,G,U,oe,rt,$t,et,kt,ae,ct,dt,lt,N,zt,yt,y,st,vt,Jt,Ft,j,nt,Lt,Kt,At,F,pt,Tt,ie,ft,re,D,Pt,ot,St,O,ut,ce,z,Ot,Qt,Rt,de,q,Vt,J,mt,pe,I,fe,B,ue,P,Et,K,ht,me,bt,he,w,Nt,at,qt,be,Ht,Wt,Q,gt,ge,Mt,_e,_t,xe,V,xt,we,wt,Xt,W,Yt,L,X,R,Dt,Ce,Y,v,It;return{c(){n=e("p"),n.innerHTML=`The syntax basically follows the format
        <code><span class="txt-success">OPERAND</span> <span class="txt-danger">OPERATOR</span> <span class="txt-success">OPERAND</span></code>, where:`,o=s(),i=e("ul"),f=e("li"),f.innerHTML=`<code class="txt-success">OPERAND</code> - could be any of the above field literal, string (single
            or double quoted), number, null, true, false`,h=s(),r=e("li"),b=e("code"),b.textContent="OPERATOR",C=_(` - is one of:
            `),$=e("br"),g=s(),p=e("ul"),Z=e("li"),Ct=e("code"),Ct.textContent="=",Ut=s(),E=e("span"),E.textContent="Equal",jt=s(),M=e("li"),it=e("code"),it.textContent="!=",S=s(),tt=e("span"),tt.textContent="NOT equal",ne=s(),G=e("li"),U=e("code"),U.textContent=">",oe=s(),rt=e("span"),rt.textContent="Greater than",$t=s(),et=e("li"),kt=e("code"),kt.textContent=">=",ae=s(),ct=e("span"),ct.textContent="Greater than or equal",dt=s(),lt=e("li"),N=e("code"),N.textContent="<",zt=s(),yt=e("span"),yt.textContent="Less than",y=s(),st=e("li"),vt=e("code"),vt.textContent="<=",Jt=s(),Ft=e("span"),Ft.textContent="Less than or equal",j=s(),nt=e("li"),Lt=e("code"),Lt.textContent="~",Kt=s(),At=e("span"),At.textContent=`Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,F=s(),pt=e("li"),Tt=e("code"),Tt.textContent="!~",ie=s(),ft=e("span"),ft.textContent=`NOT Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,re=s(),D=e("li"),Pt=e("code"),Pt.textContent="?=",ot=s(),St=e("em"),St.textContent="Any/At least one of",O=s(),ut=e("span"),ut.textContent="Equal",ce=s(),z=e("li"),Ot=e("code"),Ot.textContent="?!=",Qt=s(),Rt=e("em"),Rt.textContent="Any/At least one of",de=s(),q=e("span"),q.textContent="NOT equal",Vt=s(),J=e("li"),mt=e("code"),mt.textContent="?>",pe=s(),I=e("em"),I.textContent="Any/At least one of",fe=s(),B=e("span"),B.textContent="Greater than",ue=s(),P=e("li"),Et=e("code"),Et.textContent="?>=",K=s(),ht=e("em"),ht.textContent="Any/At least one of",me=s(),bt=e("span"),bt.textContent="Greater than or equal",he=s(),w=e("li"),Nt=e("code"),Nt.textContent="?<",at=s(),qt=e("em"),qt.textContent="Any/At least one of",be=s(),Ht=e("span"),Ht.textContent="Less than",Wt=s(),Q=e("li"),gt=e("code"),gt.textContent="?<=",ge=s(),Mt=e("em"),Mt.textContent="Any/At least one of",_e=s(),_t=e("span"),_t.textContent="Less than or equal",xe=s(),V=e("li"),xt=e("code"),xt.textContent="?~",we=s(),wt=e("em"),wt.textContent="Any/At least one of",Xt=s(),W=e("span"),W.textContent=`Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,Yt=s(),L=e("li"),X=e("code"),X.textContent="?!~",R=s(),Dt=e("em"),Dt.textContent="Any/At least one of",Ce=s(),Y=e("span"),Y.textContent=`NOT Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,v=s(),It=e("p"),It.innerHTML=`To group and combine several expressions you could use brackets
        <code>(...)</code>, <code>&amp;&amp;</code> (AND) and <code>||</code> (OR) tokens.`,a(b,"class","txt-danger"),a(Ct,"class","filter-op svelte-1w7s5nw"),a(E,"class","txt"),a(it,"class","filter-op svelte-1w7s5nw"),a(tt,"class","txt"),a(U,"class","filter-op svelte-1w7s5nw"),a(rt,"class","txt"),a(kt,"class","filter-op svelte-1w7s5nw"),a(ct,"class","txt"),a(N,"class","filter-op svelte-1w7s5nw"),a(yt,"class","txt"),a(vt,"class","filter-op svelte-1w7s5nw"),a(Ft,"class","txt"),a(Lt,"class","filter-op svelte-1w7s5nw"),a(At,"class","txt"),a(Tt,"class","filter-op svelte-1w7s5nw"),a(ft,"class","txt"),a(Pt,"class","filter-op svelte-1w7s5nw"),a(St,"class","txt-hint"),a(ut,"class","txt"),a(Ot,"class","filter-op svelte-1w7s5nw"),a(Rt,"class","txt-hint"),a(q,"class","txt"),a(mt,"class","filter-op svelte-1w7s5nw"),a(I,"class","txt-hint"),a(B,"class","txt"),a(Et,"class","filter-op svelte-1w7s5nw"),a(ht,"class","txt-hint"),a(bt,"class","txt"),a(Nt,"class","filter-op svelte-1w7s5nw"),a(qt,"class","txt-hint"),a(Ht,"class","txt"),a(gt,"class","filter-op svelte-1w7s5nw"),a(Mt,"class","txt-hint"),a(_t,"class","txt"),a(xt,"class","filter-op svelte-1w7s5nw"),a(wt,"class","txt-hint"),a(W,"class","txt"),a(X,"class","filter-op svelte-1w7s5nw"),a(Dt,"class","txt-hint"),a(Y,"class","txt")},m(A,k){u(A,n,k),u(A,o,k),u(A,i,k),t(i,f),t(i,h),t(i,r),t(r,b),t(r,C),t(r,$),t(r,g),t(r,p),t(p,Z),t(Z,Ct),t(Z,Ut),t(Z,E),t(p,jt),t(p,M),t(M,it),t(M,S),t(M,tt),t(p,ne),t(p,G),t(G,U),t(G,oe),t(G,rt),t(p,$t),t(p,et),t(et,kt),t(et,ae),t(et,ct),t(p,dt),t(p,lt),t(lt,N),t(lt,zt),t(lt,yt),t(p,y),t(p,st),t(st,vt),t(st,Jt),t(st,Ft),t(p,j),t(p,nt),t(nt,Lt),t(nt,Kt),t(nt,At),t(p,F),t(p,pt),t(pt,Tt),t(pt,ie),t(pt,ft),t(p,re),t(p,D),t(D,Pt),t(D,ot),t(D,St),t(D,O),t(D,ut),t(p,ce),t(p,z),t(z,Ot),t(z,Qt),t(z,Rt),t(z,de),t(z,q),t(p,Vt),t(p,J),t(J,mt),t(J,pe),t(J,I),t(J,fe),t(J,B),t(p,ue),t(p,P),t(P,Et),t(P,K),t(P,ht),t(P,me),t(P,bt),t(p,he),t(p,w),t(w,Nt),t(w,at),t(w,qt),t(w,be),t(w,Ht),t(p,Wt),t(p,Q),t(Q,gt),t(Q,ge),t(Q,Mt),t(Q,_e),t(Q,_t),t(p,xe),t(p,V),t(V,xt),t(V,we),t(V,wt),t(V,Xt),t(V,W),t(p,Yt),t(p,L),t(L,X),t(L,R),t(L,Dt),t(L,Ce),t(L,Y),u(A,v,k),u(A,It,k)},d(A){A&&(m(n),m(o),m(i),m(v),m(It))}}}function ul(d){let n,o,i,f,h;function r(g,p){return g[0]?fl:pl}let b=r(d),C=b(d),$=d[0]&&ze();return{c(){n=e("button"),C.c(),o=s(),$&&$.c(),i=sl(),a(n,"class","btn btn-sm btn-secondary m-t-10")},m(g,p){u(g,n,p),C.m(n,null),u(g,o,p),$&&$.m(g,p),u(g,i,p),f||(h=ll(n,"click",d[1]),f=!0)},p(g,[p]){b!==(b=r(g))&&(C.d(1),C=b(g),C&&(C.c(),C.m(n,null))),g[0]?$||($=ze(),$.c(),$.m(i.parentNode,i)):$&&($.d(1),$=null)},i:Ue,o:Ue,d(g){g&&(m(n),m(o),m(i)),C.d(),$&&$.d(g),f=!1,h()}}}function ml(d,n,o){let i=!1;function f(){o(0,i=!i)}return[i,f]}class hl extends Ze{constructor(n){super(),tl(this,n,ml,ul,el,{})}}function Je(d,n,o){const i=d.slice();return i[7]=n[o],i}function Ke(d,n,o){const i=d.slice();return i[7]=n[o],i}function Qe(d,n,o){const i=d.slice();return i[12]=n[o],i[14]=o,i}function Ve(d){let n;return{c(){n=e("p"),n.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",a(n,"class","txt-hint txt-sm txt-right")},m(o,i){u(o,n,i)},d(o){o&&m(n)}}}function We(d){let n,o=d[12]+"",i,f=d[14]<d[4].length-1?", ":"",h;return{c(){n=e("code"),i=_(o),h=_(f)},m(r,b){u(r,n,b),t(n,i),u(r,h,b)},p(r,b){b&16&&o!==(o=r[12]+"")&&ke(i,o),b&16&&f!==(f=r[14]<r[4].length-1?", ":"")&&ke(h,f)},d(r){r&&(m(n),m(h))}}}function Xe(d,n){let o,i,f;function h(){return n[6](n[7])}return{key:d,first:null,c(){o=e("button"),o.textContent=`${n[7].code} `,a(o,"type","button"),a(o,"class","tab-item"),Le(o,"active",n[2]===n[7].code),this.first=o},m(r,b){u(r,o,b),i||(f=ll(o,"click",h),i=!0)},p(r,b){n=r,b&36&&Le(o,"active",n[2]===n[7].code)},d(r){r&&m(o),i=!1,f()}}}function Ye(d,n){let o,i,f,h;return i=new Fe({props:{content:n[7].body}}),{key:d,first:null,c(){o=e("div"),ee(i.$$.fragment),f=s(),a(o,"class","tab-item"),Le(o,"active",n[2]===n[7].code),this.first=o},m(r,b){u(r,o,b),le(i,o,null),t(o,f),h=!0},p(r,b){n=r,(!h||b&36)&&Le(o,"active",n[2]===n[7].code)},i(r){h||(Bt(i.$$.fragment,r),h=!0)},o(r){Gt(i.$$.fragment,r),h=!1},d(r){r&&m(o),se(i)}}}function bl(d){var Pe,Se,Oe,Re,Ee,Ne;let n,o,i=d[0].name+"",f,h,r,b,C,$,g,p=d[0].name+"",Z,Ct,Ut,E,jt,M,it,S,tt,ne,G,U,oe,rt,$t=d[0].name+"",et,kt,ae,ct,dt,lt,N,zt,yt,y,st,vt,Jt,Ft,j,nt,Lt,Kt,At,F,pt,Tt,ie,ft,re,D,Pt,ot,St,O,ut,ce,z,Ot,Qt,Rt,de,q,Vt,J,mt,pe,I,fe,B,ue,P,Et,K,ht,me,bt,he,w,Nt,at,qt,be,Ht,Wt,Q,gt,ge,Mt,_e,_t,xe,V,xt,we,wt,Xt,W,Yt,L,X,R=[],Dt=new Map,Ce,Y,v=[],It=new Map,A;E=new dl({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${d[3]}');

        ...

        // fetch a paginated records list
        const resultList = await pb.collection('${(Pe=d[0])==null?void 0:Pe.name}').getList(1, 50, {
            filter: 'created >= "2022-01-01 00:00:00" && someField1 != someField2',
        });

        // you can also fetch all records at once via getFullList
        const records = await pb.collection('${(Se=d[0])==null?void 0:Se.name}').getFullList({
            sort: '-created',
        });

        // or fetch only the first record that matches the specified filter
        const record = await pb.collection('${(Oe=d[0])==null?void 0:Oe.name}').getFirstListItem('someField="test"', {
            expand: 'relField1,relField2.subRelField',
        });
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${d[3]}');

        ...

        // fetch a paginated records list
        final resultList = await pb.collection('${(Re=d[0])==null?void 0:Re.name}').getList(
          page: 1,
          perPage: 50,
          filter: 'created >= "2022-01-01 00:00:00" && someField1 != someField2',
        );

        // you can also fetch all records at once via getFullList
        final records = await pb.collection('${(Ee=d[0])==null?void 0:Ee.name}').getFullList(
          sort: '-created',
        );

        // or fetch only the first record that matches the specified filter
        final record = await pb.collection('${(Ne=d[0])==null?void 0:Ne.name}').getFirstListItem(
          'someField="test"',
          expand: 'relField1,relField2.subRelField',
        );
    `}});let k=d[1]&&Ve();ot=new Fe({props:{content:`
                        // DESC by created and ASC by id
                        ?sort=-created,id
                    `}});let Zt=te(d[4]),T=[];for(let l=0;l<Zt.length;l+=1)T[l]=We(Qe(d,Zt,l));B=new Fe({props:{content:`
                        ?filter=(id='abc' && created>'2022-01-01')
                    `}}),P=new hl({}),at=new Fe({props:{content:"?expand=relField1,relField2.subRelField"}});let ye=te(d[5]);const Ae=l=>l[7].code;for(let l=0;l<ye.length;l+=1){let c=Ke(d,ye,l),x=Ae(c);Dt.set(x,R[l]=Xe(x,c))}let $e=te(d[5]);const Te=l=>l[7].code;for(let l=0;l<$e.length;l+=1){let c=Je(d,$e,l),x=Te(c);It.set(x,v[l]=Ye(x,c))}return{c(){n=e("h3"),o=_("List/Search ("),f=_(i),h=_(")"),r=s(),b=e("div"),C=e("p"),$=_("Fetch a paginated "),g=e("strong"),Z=_(p),Ct=_(" records list, supporting sorting and filtering."),Ut=s(),ee(E.$$.fragment),jt=s(),M=e("h6"),M.textContent="API details",it=s(),S=e("div"),tt=e("strong"),tt.textContent="GET",ne=s(),G=e("div"),U=e("p"),oe=_("/api/collections/"),rt=e("strong"),et=_($t),kt=_("/records"),ae=s(),k&&k.c(),ct=s(),dt=e("div"),dt.textContent="Query parameters",lt=s(),N=e("table"),zt=e("thead"),zt.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',yt=s(),y=e("tbody"),st=e("tr"),st.innerHTML='<td>page</td> <td><span class="label">Number</span></td> <td>The page (aka. offset) of the paginated list (default to 1).</td>',vt=s(),Jt=e("tr"),Jt.innerHTML='<td>perPage</td> <td><span class="label">Number</span></td> <td>Specify the max returned records per page (default to 30).</td>',Ft=s(),j=e("tr"),nt=e("td"),nt.textContent="sort",Lt=s(),Kt=e("td"),Kt.innerHTML='<span class="label">String</span>',At=s(),F=e("td"),pt=_("Specify the records order attribute(s). "),Tt=e("br"),ie=_(`
                Add `),ft=e("code"),ft.textContent="-",re=_(" / "),D=e("code"),D.textContent="+",Pt=_(` (default) in front of the attribute for DESC / ASC order.
                Ex.:
                `),ee(ot.$$.fragment),St=s(),O=e("p"),ut=e("strong"),ut.textContent="Supported record sort fields:",ce=s(),z=e("br"),Ot=s(),Qt=e("code"),Qt.textContent="@random",Rt=_(`,
                    `);for(let l=0;l<T.length;l+=1)T[l].c();de=s(),q=e("tr"),Vt=e("td"),Vt.textContent="filter",J=s(),mt=e("td"),mt.innerHTML='<span class="label">String</span>',pe=s(),I=e("td"),fe=_(`Filter the returned records. Ex.:
                `),ee(B.$$.fragment),ue=s(),ee(P.$$.fragment),Et=s(),K=e("tr"),ht=e("td"),ht.textContent="expand",me=s(),bt=e("td"),bt.innerHTML='<span class="label">String</span>',he=s(),w=e("td"),Nt=_(`Auto expand record relations. Ex.:
                `),ee(at.$$.fragment),qt=_(`
                Supports up to 6-levels depth nested relations expansion. `),be=e("br"),Ht=_(`
                The expanded relations will be appended to each individual record under the
                `),Wt=e("code"),Wt.textContent="expand",Q=_(" property (eg. "),gt=e("code"),gt.textContent='"expand": {"relField1": {...}, ...}',ge=_(`).
                `),Mt=e("br"),_e=_(`
                Only the relations to which the request user has permissions to `),_t=e("strong"),_t.textContent="view",xe=_(" will be expanded."),V=s(),xt=e("tr"),xt.innerHTML=`<td id="query-page">fields</td> <td><span class="label">String</span></td> <td>Comma separated string of the fields to return in the JSON response
                <em>(by default returns all fields)</em>.
                <br/>
                Use <code>*</code> to target all keys from the specific depth level.
                <br/>
                For example: <code>?fields=*,expand.author.name</code></td>`,we=s(),wt=e("tr"),wt.innerHTML=`<td id="query-page">skipTotal</td> <td><span class="label">Boolean</span></td> <td>If it is set the total counts query will be skipped and the response fields
                <code>totalItems</code> and <code>totalPages</code> will have <code>-1</code> value.
                <br/>
                This could drastically speed up the search queries when the total counters are not needed or cursor
                based pagination is used.
                <br/>
                For optimization purposes, it is set by default for the
                <code>getFirstListItem()</code>
                and
                <code>getFullList()</code> SDKs methods.</td>`,Xt=s(),W=e("div"),W.textContent="Responses",Yt=s(),L=e("div"),X=e("div");for(let l=0;l<R.length;l+=1)R[l].c();Ce=s(),Y=e("div");for(let l=0;l<v.length;l+=1)v[l].c();a(n,"class","m-b-sm"),a(b,"class","content txt-lg m-b-sm"),a(M,"class","m-b-xs"),a(tt,"class","label label-primary"),a(G,"class","content"),a(S,"class","alert alert-info"),a(dt,"class","section-title"),a(N,"class","table-compact table-border m-b-base"),a(W,"class","section-title"),a(X,"class","tabs-header compact combined left"),a(Y,"class","tabs-content"),a(L,"class","tabs")},m(l,c){u(l,n,c),t(n,o),t(n,f),t(n,h),u(l,r,c),u(l,b,c),t(b,C),t(C,$),t(C,g),t(g,Z),t(C,Ct),u(l,Ut,c),le(E,l,c),u(l,jt,c),u(l,M,c),u(l,it,c),u(l,S,c),t(S,tt),t(S,ne),t(S,G),t(G,U),t(U,oe),t(U,rt),t(rt,et),t(U,kt),t(S,ae),k&&k.m(S,null),u(l,ct,c),u(l,dt,c),u(l,lt,c),u(l,N,c),t(N,zt),t(N,yt),t(N,y),t(y,st),t(y,vt),t(y,Jt),t(y,Ft),t(y,j),t(j,nt),t(j,Lt),t(j,Kt),t(j,At),t(j,F),t(F,pt),t(F,Tt),t(F,ie),t(F,ft),t(F,re),t(F,D),t(F,Pt),le(ot,F,null),t(F,St),t(F,O),t(O,ut),t(O,ce),t(O,z),t(O,Ot),t(O,Qt),t(O,Rt);for(let x=0;x<T.length;x+=1)T[x]&&T[x].m(O,null);t(y,de),t(y,q),t(q,Vt),t(q,J),t(q,mt),t(q,pe),t(q,I),t(I,fe),le(B,I,null),t(I,ue),le(P,I,null),t(y,Et),t(y,K),t(K,ht),t(K,me),t(K,bt),t(K,he),t(K,w),t(w,Nt),le(at,w,null),t(w,qt),t(w,be),t(w,Ht),t(w,Wt),t(w,Q),t(w,gt),t(w,ge),t(w,Mt),t(w,_e),t(w,_t),t(w,xe),t(y,V),t(y,xt),t(y,we),t(y,wt),u(l,Xt,c),u(l,W,c),u(l,Yt,c),u(l,L,c),t(L,X);for(let x=0;x<R.length;x+=1)R[x]&&R[x].m(X,null);t(L,Ce),t(L,Y);for(let x=0;x<v.length;x+=1)v[x]&&v[x].m(Y,null);A=!0},p(l,[c]){var qe,He,Me,De,Ie,Be;(!A||c&1)&&i!==(i=l[0].name+"")&&ke(f,i),(!A||c&1)&&p!==(p=l[0].name+"")&&ke(Z,p);const x={};if(c&9&&(x.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${l[3]}');

        ...

        // fetch a paginated records list
        const resultList = await pb.collection('${(qe=l[0])==null?void 0:qe.name}').getList(1, 50, {
            filter: 'created >= "2022-01-01 00:00:00" && someField1 != someField2',
        });

        // you can also fetch all records at once via getFullList
        const records = await pb.collection('${(He=l[0])==null?void 0:He.name}').getFullList({
            sort: '-created',
        });

        // or fetch only the first record that matches the specified filter
        const record = await pb.collection('${(Me=l[0])==null?void 0:Me.name}').getFirstListItem('someField="test"', {
            expand: 'relField1,relField2.subRelField',
        });
    `),c&9&&(x.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${l[3]}');

        ...

        // fetch a paginated records list
        final resultList = await pb.collection('${(De=l[0])==null?void 0:De.name}').getList(
          page: 1,
          perPage: 50,
          filter: 'created >= "2022-01-01 00:00:00" && someField1 != someField2',
        );

        // you can also fetch all records at once via getFullList
        final records = await pb.collection('${(Ie=l[0])==null?void 0:Ie.name}').getFullList(
          sort: '-created',
        );

        // or fetch only the first record that matches the specified filter
        final record = await pb.collection('${(Be=l[0])==null?void 0:Be.name}').getFirstListItem(
          'someField="test"',
          expand: 'relField1,relField2.subRelField',
        );
    `),E.$set(x),(!A||c&1)&&$t!==($t=l[0].name+"")&&ke(et,$t),l[1]?k||(k=Ve(),k.c(),k.m(S,null)):k&&(k.d(1),k=null),c&16){Zt=te(l[4]);let H;for(H=0;H<Zt.length;H+=1){const Ge=Qe(l,Zt,H);T[H]?T[H].p(Ge,c):(T[H]=We(Ge),T[H].c(),T[H].m(O,null))}for(;H<T.length;H+=1)T[H].d(1);T.length=Zt.length}c&36&&(ye=te(l[5]),R=je(R,c,Ae,1,l,ye,Dt,X,nl,Xe,null,Ke)),c&36&&($e=te(l[5]),ol(),v=je(v,c,Te,1,l,$e,It,Y,al,Ye,null,Je),il())},i(l){if(!A){Bt(E.$$.fragment,l),Bt(ot.$$.fragment,l),Bt(B.$$.fragment,l),Bt(P.$$.fragment,l),Bt(at.$$.fragment,l);for(let c=0;c<$e.length;c+=1)Bt(v[c]);A=!0}},o(l){Gt(E.$$.fragment,l),Gt(ot.$$.fragment,l),Gt(B.$$.fragment,l),Gt(P.$$.fragment,l),Gt(at.$$.fragment,l);for(let c=0;c<v.length;c+=1)Gt(v[c]);A=!1},d(l){l&&(m(n),m(r),m(b),m(Ut),m(jt),m(M),m(it),m(S),m(ct),m(dt),m(lt),m(N),m(Xt),m(W),m(Yt),m(L)),se(E,l),k&&k.d(),se(ot),rl(T,l),se(B),se(P),se(at);for(let c=0;c<R.length;c+=1)R[c].d();for(let c=0;c<v.length;c+=1)v[c].d()}}}function gl(d,n,o){let i,f,h,{collection:r}=n,b=200,C=[];const $=g=>o(2,b=g.code);return d.$$set=g=>{"collection"in g&&o(0,r=g.collection)},d.$$.update=()=>{d.$$.dirty&1&&o(4,i=ve.getAllCollectionIdentifiers(r)),d.$$.dirty&1&&o(1,f=(r==null?void 0:r.listRule)===null),d.$$.dirty&3&&r!=null&&r.id&&(C.push({code:200,body:JSON.stringify({page:1,perPage:30,totalPages:1,totalItems:2,items:[ve.dummyCollectionRecord(r),ve.dummyCollectionRecord(r)]},null,2)}),C.push({code:400,body:`
                {
                  "code": 400,
                  "message": "Something went wrong while processing your request. Invalid filter.",
                  "data": {}
                }
            `}),f&&C.push({code:403,body:`
                    {
                      "code": 403,
                      "message": "Only admins can access this action.",
                      "data": {}
                    }
                `}))},o(3,h=ve.getApiExampleUrl(cl.baseUrl)),[r,f,b,h,i,C,$]}class wl extends Ze{constructor(n){super(),tl(this,n,gl,bl,el,{collection:0})}}export{wl as default};
