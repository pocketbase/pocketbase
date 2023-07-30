import{S as Ye,i as Ze,s as tl,e,b as s,E as ll,f as i,g as u,u as el,y as Ge,o as m,w as x,h as t,M as ve,c as te,m as ee,x as $e,N as Ue,P as sl,k as nl,Q as ol,n as il,t as Bt,a as Gt,d as le,R as al,T as rl,C as ye,p as cl,r as Fe}from"./index-c8d873a4.js";import{S as dl}from"./SdkTabs-2977cee7.js";function fl(c){let n,o,a;return{c(){n=e("span"),n.textContent="Show details",o=s(),a=e("i"),i(n,"class","txt"),i(a,"class","ri-arrow-down-s-line")},m(p,b){u(p,n,b),u(p,o,b),u(p,a,b)},d(p){p&&m(n),p&&m(o),p&&m(a)}}}function pl(c){let n,o,a;return{c(){n=e("span"),n.textContent="Hide details",o=s(),a=e("i"),i(n,"class","txt"),i(a,"class","ri-arrow-up-s-line")},m(p,b){u(p,n,b),u(p,o,b),u(p,a,b)},d(p){p&&m(n),p&&m(o),p&&m(a)}}}function je(c){let n,o,a,p,b,d,h,g,w,_,f,Z,Ct,Ut,E,jt,H,at,S,tt,se,G,U,ne,rt,$t,et,kt,oe,ct,dt,lt,N,zt,yt,v,st,vt,Jt,Ft,j,nt,Lt,Kt,At,L,ft,Tt,ie,pt,ae,D,Pt,ot,St,R,ut,re,z,Rt,Qt,Ot,ce,q,Vt,J,mt,de,I,fe,B,pe,P,Et,K,bt,ue,ht,me,$,Nt,it,qt,be,Mt,Wt,Q,_t,he,Ht,_e,wt,we,V,xt,xe,gt,Xt,W,Yt,A,X,O,Dt,ge,Y,F,It;return{c(){n=e("p"),n.innerHTML=`The syntax basically follows the format
        <code><span class="txt-success">OPERAND</span> 
            <span class="txt-danger">OPERATOR</span> 
            <span class="txt-success">OPERAND</span></code>, where:`,o=s(),a=e("ul"),p=e("li"),p.innerHTML=`<code class="txt-success">OPERAND</code> - could be any of the above field literal, string (single
            or double quoted), number, null, true, false`,b=s(),d=e("li"),h=e("code"),h.textContent="OPERATOR",g=x(` - is one of:
            `),w=e("br"),_=s(),f=e("ul"),Z=e("li"),Ct=e("code"),Ct.textContent="=",Ut=s(),E=e("span"),E.textContent="Equal",jt=s(),H=e("li"),at=e("code"),at.textContent="!=",S=s(),tt=e("span"),tt.textContent="NOT equal",se=s(),G=e("li"),U=e("code"),U.textContent=">",ne=s(),rt=e("span"),rt.textContent="Greater than",$t=s(),et=e("li"),kt=e("code"),kt.textContent=">=",oe=s(),ct=e("span"),ct.textContent="Greater than or equal",dt=s(),lt=e("li"),N=e("code"),N.textContent="<",zt=s(),yt=e("span"),yt.textContent="Less than",v=s(),st=e("li"),vt=e("code"),vt.textContent="<=",Jt=s(),Ft=e("span"),Ft.textContent="Less than or equal",j=s(),nt=e("li"),Lt=e("code"),Lt.textContent="~",Kt=s(),At=e("span"),At.textContent=`Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,L=s(),ft=e("li"),Tt=e("code"),Tt.textContent="!~",ie=s(),pt=e("span"),pt.textContent=`NOT Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,ae=s(),D=e("li"),Pt=e("code"),Pt.textContent="?=",ot=s(),St=e("em"),St.textContent="Any/At least one of",R=s(),ut=e("span"),ut.textContent="Equal",re=s(),z=e("li"),Rt=e("code"),Rt.textContent="?!=",Qt=s(),Ot=e("em"),Ot.textContent="Any/At least one of",ce=s(),q=e("span"),q.textContent="NOT equal",Vt=s(),J=e("li"),mt=e("code"),mt.textContent="?>",de=s(),I=e("em"),I.textContent="Any/At least one of",fe=s(),B=e("span"),B.textContent="Greater than",pe=s(),P=e("li"),Et=e("code"),Et.textContent="?>=",K=s(),bt=e("em"),bt.textContent="Any/At least one of",ue=s(),ht=e("span"),ht.textContent="Greater than or equal",me=s(),$=e("li"),Nt=e("code"),Nt.textContent="?<",it=s(),qt=e("em"),qt.textContent="Any/At least one of",be=s(),Mt=e("span"),Mt.textContent="Less than",Wt=s(),Q=e("li"),_t=e("code"),_t.textContent="?<=",he=s(),Ht=e("em"),Ht.textContent="Any/At least one of",_e=s(),wt=e("span"),wt.textContent="Less than or equal",we=s(),V=e("li"),xt=e("code"),xt.textContent="?~",xe=s(),gt=e("em"),gt.textContent="Any/At least one of",Xt=s(),W=e("span"),W.textContent=`Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,Yt=s(),A=e("li"),X=e("code"),X.textContent="?!~",O=s(),Dt=e("em"),Dt.textContent="Any/At least one of",ge=s(),Y=e("span"),Y.textContent=`NOT Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,F=s(),It=e("p"),It.innerHTML=`To group and combine several expressions you could use brackets
        <code>(...)</code>, <code>&amp;&amp;</code> (AND) and <code>||</code> (OR) tokens.`,i(h,"class","txt-danger"),i(Ct,"class","filter-op svelte-1w7s5nw"),i(E,"class","txt"),i(at,"class","filter-op svelte-1w7s5nw"),i(tt,"class","txt"),i(U,"class","filter-op svelte-1w7s5nw"),i(rt,"class","txt"),i(kt,"class","filter-op svelte-1w7s5nw"),i(ct,"class","txt"),i(N,"class","filter-op svelte-1w7s5nw"),i(yt,"class","txt"),i(vt,"class","filter-op svelte-1w7s5nw"),i(Ft,"class","txt"),i(Lt,"class","filter-op svelte-1w7s5nw"),i(At,"class","txt"),i(Tt,"class","filter-op svelte-1w7s5nw"),i(pt,"class","txt"),i(Pt,"class","filter-op svelte-1w7s5nw"),i(St,"class","txt-hint"),i(ut,"class","txt"),i(Rt,"class","filter-op svelte-1w7s5nw"),i(Ot,"class","txt-hint"),i(q,"class","txt"),i(mt,"class","filter-op svelte-1w7s5nw"),i(I,"class","txt-hint"),i(B,"class","txt"),i(Et,"class","filter-op svelte-1w7s5nw"),i(bt,"class","txt-hint"),i(ht,"class","txt"),i(Nt,"class","filter-op svelte-1w7s5nw"),i(qt,"class","txt-hint"),i(Mt,"class","txt"),i(_t,"class","filter-op svelte-1w7s5nw"),i(Ht,"class","txt-hint"),i(wt,"class","txt"),i(xt,"class","filter-op svelte-1w7s5nw"),i(gt,"class","txt-hint"),i(W,"class","txt"),i(X,"class","filter-op svelte-1w7s5nw"),i(Dt,"class","txt-hint"),i(Y,"class","txt")},m(k,y){u(k,n,y),u(k,o,y),u(k,a,y),t(a,p),t(a,b),t(a,d),t(d,h),t(d,g),t(d,w),t(d,_),t(d,f),t(f,Z),t(Z,Ct),t(Z,Ut),t(Z,E),t(f,jt),t(f,H),t(H,at),t(H,S),t(H,tt),t(f,se),t(f,G),t(G,U),t(G,ne),t(G,rt),t(f,$t),t(f,et),t(et,kt),t(et,oe),t(et,ct),t(f,dt),t(f,lt),t(lt,N),t(lt,zt),t(lt,yt),t(f,v),t(f,st),t(st,vt),t(st,Jt),t(st,Ft),t(f,j),t(f,nt),t(nt,Lt),t(nt,Kt),t(nt,At),t(f,L),t(f,ft),t(ft,Tt),t(ft,ie),t(ft,pt),t(f,ae),t(f,D),t(D,Pt),t(D,ot),t(D,St),t(D,R),t(D,ut),t(f,re),t(f,z),t(z,Rt),t(z,Qt),t(z,Ot),t(z,ce),t(z,q),t(f,Vt),t(f,J),t(J,mt),t(J,de),t(J,I),t(J,fe),t(J,B),t(f,pe),t(f,P),t(P,Et),t(P,K),t(P,bt),t(P,ue),t(P,ht),t(f,me),t(f,$),t($,Nt),t($,it),t($,qt),t($,be),t($,Mt),t(f,Wt),t(f,Q),t(Q,_t),t(Q,he),t(Q,Ht),t(Q,_e),t(Q,wt),t(f,we),t(f,V),t(V,xt),t(V,xe),t(V,gt),t(V,Xt),t(V,W),t(f,Yt),t(f,A),t(A,X),t(A,O),t(A,Dt),t(A,ge),t(A,Y),u(k,F,y),u(k,It,y)},d(k){k&&m(n),k&&m(o),k&&m(a),k&&m(F),k&&m(It)}}}function ul(c){let n,o,a,p,b;function d(_,f){return _[0]?pl:fl}let h=d(c),g=h(c),w=c[0]&&je();return{c(){n=e("button"),g.c(),o=s(),w&&w.c(),a=ll(),i(n,"class","btn btn-sm btn-secondary m-t-10")},m(_,f){u(_,n,f),g.m(n,null),u(_,o,f),w&&w.m(_,f),u(_,a,f),p||(b=el(n,"click",c[1]),p=!0)},p(_,[f]){h!==(h=d(_))&&(g.d(1),g=h(_),g&&(g.c(),g.m(n,null))),_[0]?w||(w=je(),w.c(),w.m(a.parentNode,a)):w&&(w.d(1),w=null)},i:Ge,o:Ge,d(_){_&&m(n),g.d(),_&&m(o),w&&w.d(_),_&&m(a),p=!1,b()}}}function ml(c,n,o){let a=!1;function p(){o(0,a=!a)}return[a,p]}class bl extends Ye{constructor(n){super(),Ze(this,n,ml,ul,tl,{})}}function ze(c,n,o){const a=c.slice();return a[7]=n[o],a}function Je(c,n,o){const a=c.slice();return a[7]=n[o],a}function Ke(c,n,o){const a=c.slice();return a[12]=n[o],a[14]=o,a}function Qe(c){let n;return{c(){n=e("p"),n.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",i(n,"class","txt-hint txt-sm txt-right")},m(o,a){u(o,n,a)},d(o){o&&m(n)}}}function Ve(c){let n,o=c[12]+"",a,p=c[14]<c[4].length-1?", ":"",b;return{c(){n=e("code"),a=x(o),b=x(p)},m(d,h){u(d,n,h),t(n,a),u(d,b,h)},p(d,h){h&16&&o!==(o=d[12]+"")&&$e(a,o),h&16&&p!==(p=d[14]<d[4].length-1?", ":"")&&$e(b,p)},d(d){d&&m(n),d&&m(b)}}}function We(c,n){let o,a=n[7].code+"",p,b,d,h;function g(){return n[6](n[7])}return{key:c,first:null,c(){o=e("button"),p=x(a),b=s(),i(o,"type","button"),i(o,"class","tab-item"),Fe(o,"active",n[2]===n[7].code),this.first=o},m(w,_){u(w,o,_),t(o,p),t(o,b),d||(h=el(o,"click",g),d=!0)},p(w,_){n=w,_&36&&Fe(o,"active",n[2]===n[7].code)},d(w){w&&m(o),d=!1,h()}}}function Xe(c,n){let o,a,p,b;return a=new ve({props:{content:n[7].body}}),{key:c,first:null,c(){o=e("div"),te(a.$$.fragment),p=s(),i(o,"class","tab-item"),Fe(o,"active",n[2]===n[7].code),this.first=o},m(d,h){u(d,o,h),ee(a,o,null),t(o,p),b=!0},p(d,h){n=d,(!b||h&36)&&Fe(o,"active",n[2]===n[7].code)},i(d){b||(Bt(a.$$.fragment,d),b=!0)},o(d){Gt(a.$$.fragment,d),b=!1},d(d){d&&m(o),le(a)}}}function hl(c){var Te,Pe,Se,Re,Oe,Ee;let n,o,a=c[0].name+"",p,b,d,h,g,w,_,f=c[0].name+"",Z,Ct,Ut,E,jt,H,at,S,tt,se,G,U,ne,rt,$t=c[0].name+"",et,kt,oe,ct,dt,lt,N,zt,yt,v,st,vt,Jt,Ft,j,nt,Lt,Kt,At,L,ft,Tt,ie,pt,ae,D,Pt,ot,St,R,ut,re,z,Rt,Qt,Ot,ce,q,Vt,J,mt,de,I,fe,B,pe,P,Et,K,bt,ue,ht,me,$,Nt,it,qt,be,Mt,Wt,Q,_t,he,Ht,_e,wt,we,V,xt,xe,gt,Xt,W,Yt,A,X,O=[],Dt=new Map,ge,Y,F=[],It=new Map,k;E=new dl({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${c[3]}');

        ...

        // fetch a paginated records list
        const resultList = await pb.collection('${(Te=c[0])==null?void 0:Te.name}').getList(1, 50, {
            filter: 'created >= "2022-01-01 00:00:00" && someField1 != someField2',
        });

        // you can also fetch all records at once via getFullList
        const records = await pb.collection('${(Pe=c[0])==null?void 0:Pe.name}').getFullList({
            sort: '-created',
        });

        // or fetch only the first record that matches the specified filter
        const record = await pb.collection('${(Se=c[0])==null?void 0:Se.name}').getFirstListItem('someField="test"', {
            expand: 'relField1,relField2.subRelField',
        });
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${c[3]}');

        ...

        // fetch a paginated records list
        final resultList = await pb.collection('${(Re=c[0])==null?void 0:Re.name}').getList(
          page: 1,
          perPage: 50,
          filter: 'created >= "2022-01-01 00:00:00" && someField1 != someField2',
        );

        // you can also fetch all records at once via getFullList
        final records = await pb.collection('${(Oe=c[0])==null?void 0:Oe.name}').getFullList(
          sort: '-created',
        );

        // or fetch only the first record that matches the specified filter
        final record = await pb.collection('${(Ee=c[0])==null?void 0:Ee.name}').getFirstListItem(
          'someField="test"',
          expand: 'relField1,relField2.subRelField',
        );
    `}});let y=c[1]&&Qe();ot=new ve({props:{content:`
                        // DESC by created and ASC by id
                        ?sort=-created,id
                    `}});let Zt=c[4],T=[];for(let l=0;l<Zt.length;l+=1)T[l]=Ve(Ke(c,Zt,l));B=new ve({props:{content:`
                        ?filter=(id='abc' && created>'2022-01-01')
                    `}}),P=new bl({}),it=new ve({props:{content:"?expand=relField1,relField2.subRelField"}});let ke=c[5];const Le=l=>l[7].code;for(let l=0;l<ke.length;l+=1){let r=Je(c,ke,l),C=Le(r);Dt.set(C,O[l]=We(C,r))}let Ce=c[5];const Ae=l=>l[7].code;for(let l=0;l<Ce.length;l+=1){let r=ze(c,Ce,l),C=Ae(r);It.set(C,F[l]=Xe(C,r))}return{c(){n=e("h3"),o=x("List/Search ("),p=x(a),b=x(")"),d=s(),h=e("div"),g=e("p"),w=x("Fetch a paginated "),_=e("strong"),Z=x(f),Ct=x(" records list, supporting sorting and filtering."),Ut=s(),te(E.$$.fragment),jt=s(),H=e("h6"),H.textContent="API details",at=s(),S=e("div"),tt=e("strong"),tt.textContent="GET",se=s(),G=e("div"),U=e("p"),ne=x("/api/collections/"),rt=e("strong"),et=x($t),kt=x("/records"),oe=s(),y&&y.c(),ct=s(),dt=e("div"),dt.textContent="Query parameters",lt=s(),N=e("table"),zt=e("thead"),zt.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,yt=s(),v=e("tbody"),st=e("tr"),st.innerHTML=`<td>page</td> 
            <td><span class="label">Number</span></td> 
            <td>The page (aka. offset) of the paginated list (default to 1).</td>`,vt=s(),Jt=e("tr"),Jt.innerHTML=`<td>perPage</td> 
            <td><span class="label">Number</span></td> 
            <td>Specify the max returned records per page (default to 30).</td>`,Ft=s(),j=e("tr"),nt=e("td"),nt.textContent="sort",Lt=s(),Kt=e("td"),Kt.innerHTML='<span class="label">String</span>',At=s(),L=e("td"),ft=x("Specify the records order attribute(s). "),Tt=e("br"),ie=x(`
                Add `),pt=e("code"),pt.textContent="-",ae=x(" / "),D=e("code"),D.textContent="+",Pt=x(` (default) in front of the attribute for DESC / ASC order.
                Ex.:
                `),te(ot.$$.fragment),St=s(),R=e("p"),ut=e("strong"),ut.textContent="Supported record sort fields:",re=s(),z=e("br"),Rt=s(),Qt=e("code"),Qt.textContent="@random",Ot=x(`,
                    `);for(let l=0;l<T.length;l+=1)T[l].c();ce=s(),q=e("tr"),Vt=e("td"),Vt.textContent="filter",J=s(),mt=e("td"),mt.innerHTML='<span class="label">String</span>',de=s(),I=e("td"),fe=x(`Filter the returned records. Ex.:
                `),te(B.$$.fragment),pe=s(),te(P.$$.fragment),Et=s(),K=e("tr"),bt=e("td"),bt.textContent="expand",ue=s(),ht=e("td"),ht.innerHTML='<span class="label">String</span>',me=s(),$=e("td"),Nt=x(`Auto expand record relations. Ex.:
                `),te(it.$$.fragment),qt=x(`
                Supports up to 6-levels depth nested relations expansion. `),be=e("br"),Mt=x(`
                The expanded relations will be appended to each individual record under the
                `),Wt=e("code"),Wt.textContent="expand",Q=x(" property (eg. "),_t=e("code"),_t.textContent='"expand": {"relField1": {...}, ...}',he=x(`).
                `),Ht=e("br"),_e=x(`
                Only the relations to which the request user has permissions to `),wt=e("strong"),wt.textContent="view",we=x(" will be expanded."),V=s(),xt=e("tr"),xt.innerHTML=`<td id="query-page">fields</td> 
            <td><span class="label">String</span></td> 
            <td>Comma separated string of the fields to return in the JSON response
                <em>(by default returns all fields)</em>.</td>`,xe=s(),gt=e("tr"),gt.innerHTML=`<td id="query-page">skipTotal</td> 
            <td><span class="label">Boolean</span></td> 
            <td>If it is set the total counts query will be skipped and the response fields
                <code>totalItems</code> and <code>totalPages</code> will have <code>-1</code> value.
                <br/>
                This could drastically speed up the search queries when the total counters are not needed or cursor
                based pagination is used.
                <br/>
                For optimization purposes, it is set by default for the
                <code>getFirstListItem()</code>
                and
                <code>getFullList()</code> SDKs methods.</td>`,Xt=s(),W=e("div"),W.textContent="Responses",Yt=s(),A=e("div"),X=e("div");for(let l=0;l<O.length;l+=1)O[l].c();ge=s(),Y=e("div");for(let l=0;l<F.length;l+=1)F[l].c();i(n,"class","m-b-sm"),i(h,"class","content txt-lg m-b-sm"),i(H,"class","m-b-xs"),i(tt,"class","label label-primary"),i(G,"class","content"),i(S,"class","alert alert-info"),i(dt,"class","section-title"),i(N,"class","table-compact table-border m-b-base"),i(W,"class","section-title"),i(X,"class","tabs-header compact left"),i(Y,"class","tabs-content"),i(A,"class","tabs")},m(l,r){u(l,n,r),t(n,o),t(n,p),t(n,b),u(l,d,r),u(l,h,r),t(h,g),t(g,w),t(g,_),t(_,Z),t(g,Ct),u(l,Ut,r),ee(E,l,r),u(l,jt,r),u(l,H,r),u(l,at,r),u(l,S,r),t(S,tt),t(S,se),t(S,G),t(G,U),t(U,ne),t(U,rt),t(rt,et),t(U,kt),t(S,oe),y&&y.m(S,null),u(l,ct,r),u(l,dt,r),u(l,lt,r),u(l,N,r),t(N,zt),t(N,yt),t(N,v),t(v,st),t(v,vt),t(v,Jt),t(v,Ft),t(v,j),t(j,nt),t(j,Lt),t(j,Kt),t(j,At),t(j,L),t(L,ft),t(L,Tt),t(L,ie),t(L,pt),t(L,ae),t(L,D),t(L,Pt),ee(ot,L,null),t(L,St),t(L,R),t(R,ut),t(R,re),t(R,z),t(R,Rt),t(R,Qt),t(R,Ot);for(let C=0;C<T.length;C+=1)T[C]&&T[C].m(R,null);t(v,ce),t(v,q),t(q,Vt),t(q,J),t(q,mt),t(q,de),t(q,I),t(I,fe),ee(B,I,null),t(I,pe),ee(P,I,null),t(v,Et),t(v,K),t(K,bt),t(K,ue),t(K,ht),t(K,me),t(K,$),t($,Nt),ee(it,$,null),t($,qt),t($,be),t($,Mt),t($,Wt),t($,Q),t($,_t),t($,he),t($,Ht),t($,_e),t($,wt),t($,we),t(v,V),t(v,xt),t(v,xe),t(v,gt),u(l,Xt,r),u(l,W,r),u(l,Yt,r),u(l,A,r),t(A,X);for(let C=0;C<O.length;C+=1)O[C]&&O[C].m(X,null);t(A,ge),t(A,Y);for(let C=0;C<F.length;C+=1)F[C]&&F[C].m(Y,null);k=!0},p(l,[r]){var Ne,qe,Me,He,De,Ie;(!k||r&1)&&a!==(a=l[0].name+"")&&$e(p,a),(!k||r&1)&&f!==(f=l[0].name+"")&&$e(Z,f);const C={};if(r&9&&(C.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${l[3]}');

        ...

        // fetch a paginated records list
        const resultList = await pb.collection('${(Ne=l[0])==null?void 0:Ne.name}').getList(1, 50, {
            filter: 'created >= "2022-01-01 00:00:00" && someField1 != someField2',
        });

        // you can also fetch all records at once via getFullList
        const records = await pb.collection('${(qe=l[0])==null?void 0:qe.name}').getFullList({
            sort: '-created',
        });

        // or fetch only the first record that matches the specified filter
        const record = await pb.collection('${(Me=l[0])==null?void 0:Me.name}').getFirstListItem('someField="test"', {
            expand: 'relField1,relField2.subRelField',
        });
    `),r&9&&(C.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${l[3]}');

        ...

        // fetch a paginated records list
        final resultList = await pb.collection('${(He=l[0])==null?void 0:He.name}').getList(
          page: 1,
          perPage: 50,
          filter: 'created >= "2022-01-01 00:00:00" && someField1 != someField2',
        );

        // you can also fetch all records at once via getFullList
        final records = await pb.collection('${(De=l[0])==null?void 0:De.name}').getFullList(
          sort: '-created',
        );

        // or fetch only the first record that matches the specified filter
        final record = await pb.collection('${(Ie=l[0])==null?void 0:Ie.name}').getFirstListItem(
          'someField="test"',
          expand: 'relField1,relField2.subRelField',
        );
    `),E.$set(C),(!k||r&1)&&$t!==($t=l[0].name+"")&&$e(et,$t),l[1]?y||(y=Qe(),y.c(),y.m(S,null)):y&&(y.d(1),y=null),r&16){Zt=l[4];let M;for(M=0;M<Zt.length;M+=1){const Be=Ke(l,Zt,M);T[M]?T[M].p(Be,r):(T[M]=Ve(Be),T[M].c(),T[M].m(R,null))}for(;M<T.length;M+=1)T[M].d(1);T.length=Zt.length}r&36&&(ke=l[5],O=Ue(O,r,Le,1,l,ke,Dt,X,sl,We,null,Je)),r&36&&(Ce=l[5],nl(),F=Ue(F,r,Ae,1,l,Ce,It,Y,ol,Xe,null,ze),il())},i(l){if(!k){Bt(E.$$.fragment,l),Bt(ot.$$.fragment,l),Bt(B.$$.fragment,l),Bt(P.$$.fragment,l),Bt(it.$$.fragment,l);for(let r=0;r<Ce.length;r+=1)Bt(F[r]);k=!0}},o(l){Gt(E.$$.fragment,l),Gt(ot.$$.fragment,l),Gt(B.$$.fragment,l),Gt(P.$$.fragment,l),Gt(it.$$.fragment,l);for(let r=0;r<F.length;r+=1)Gt(F[r]);k=!1},d(l){l&&m(n),l&&m(d),l&&m(h),l&&m(Ut),le(E,l),l&&m(jt),l&&m(H),l&&m(at),l&&m(S),y&&y.d(),l&&m(ct),l&&m(dt),l&&m(lt),l&&m(N),le(ot),al(T,l),le(B),le(P),le(it),l&&m(Xt),l&&m(W),l&&m(Yt),l&&m(A);for(let r=0;r<O.length;r+=1)O[r].d();for(let r=0;r<F.length;r+=1)F[r].d()}}}function _l(c,n,o){let a,p,b,{collection:d=new rl}=n,h=200,g=[];const w=_=>o(2,h=_.code);return c.$$set=_=>{"collection"in _&&o(0,d=_.collection)},c.$$.update=()=>{c.$$.dirty&1&&o(4,a=ye.getAllCollectionIdentifiers(d)),c.$$.dirty&1&&o(1,p=(d==null?void 0:d.listRule)===null),c.$$.dirty&3&&d!=null&&d.id&&(g.push({code:200,body:JSON.stringify({page:1,perPage:30,totalPages:1,totalItems:2,items:[ye.dummyCollectionRecord(d),ye.dummyCollectionRecord(d)]},null,2)}),g.push({code:400,body:`
                {
                  "code": 400,
                  "message": "Something went wrong while processing your request. Invalid filter.",
                  "data": {}
                }
            `}),p&&g.push({code:403,body:`
                    {
                      "code": 403,
                      "message": "Only admins can access this action.",
                      "data": {}
                    }
                `}))},o(3,b=ye.getApiExampleUrl(cl.baseUrl)),[d,p,h,b,a,g,w]}class gl extends Ye{constructor(n){super(),Ze(this,n,_l,hl,tl,{collection:0})}}export{gl as default};
