import{S as Xe,i as Ye,s as Ze,e,b as n,E as el,f as a,g as u,u as tl,y as Be,o as m,w as _,h as t,M as ye,N as Zt,c as te,m as ee,x as Ce,P as Ge,Q as ll,k as nl,R as sl,n as ol,t as Bt,a as Gt,d as le,T as al,U as il,C as ke,p as rl,r as ve}from"./index-a084d9d7.js";import{S as cl}from"./SdkTabs-ba0ec979.js";function dl(d){let s,o,i;return{c(){s=e("span"),s.textContent="Show details",o=n(),i=e("i"),a(s,"class","txt"),a(i,"class","ri-arrow-down-s-line")},m(f,h){u(f,s,h),u(f,o,h),u(f,i,h)},d(f){f&&(m(s),m(o),m(i))}}}function pl(d){let s,o,i;return{c(){s=e("span"),s.textContent="Hide details",o=n(),i=e("i"),a(s,"class","txt"),a(i,"class","ri-arrow-up-s-line")},m(f,h){u(f,s,h),u(f,o,h),u(f,i,h)},d(f){f&&(m(s),m(o),m(i))}}}function Ue(d){let s,o,i,f,h,r,b,C,$,g,p,Z,Ct,Ut,O,jt,H,it,R,tt,ne,U,j,se,rt,$t,et,kt,oe,ct,dt,lt,E,Jt,yt,F,nt,vt,Qt,Ft,J,st,Lt,zt,At,A,pt,Tt,ae,ft,ie,D,Pt,ot,Rt,S,ut,re,Q,St,Kt,Ot,ce,N,Vt,z,mt,de,I,pe,B,fe,T,Et,K,ht,ue,bt,me,w,Nt,at,qt,he,Mt,Wt,V,gt,be,Ht,ge,_t,_e,W,xt,Xt,X,Yt,q,wt,y,Dt,xe,Y,L,It,G,v;return{c(){s=e("p"),s.innerHTML=`The syntax basically follows the format
        <code><span class="txt-success">OPERAND</span> <span class="txt-danger">OPERATOR</span> <span class="txt-success">OPERAND</span></code>, where:`,o=n(),i=e("ul"),f=e("li"),f.innerHTML=`<code class="txt-success">OPERAND</code> - could be any of the above field literal, string (single
            or double quoted), number, null, true, false`,h=n(),r=e("li"),b=e("code"),b.textContent="OPERATOR",C=_(` - is one of:
            `),$=e("br"),g=n(),p=e("ul"),Z=e("li"),Ct=e("code"),Ct.textContent="=",Ut=n(),O=e("span"),O.textContent="Equal",jt=n(),H=e("li"),it=e("code"),it.textContent="!=",R=n(),tt=e("span"),tt.textContent="NOT equal",ne=n(),U=e("li"),j=e("code"),j.textContent=">",se=n(),rt=e("span"),rt.textContent="Greater than",$t=n(),et=e("li"),kt=e("code"),kt.textContent=">=",oe=n(),ct=e("span"),ct.textContent="Greater than or equal",dt=n(),lt=e("li"),E=e("code"),E.textContent="<",Jt=n(),yt=e("span"),yt.textContent="Less than",F=n(),nt=e("li"),vt=e("code"),vt.textContent="<=",Qt=n(),Ft=e("span"),Ft.textContent="Less than or equal",J=n(),st=e("li"),Lt=e("code"),Lt.textContent="~",zt=n(),At=e("span"),At.textContent=`Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,A=n(),pt=e("li"),Tt=e("code"),Tt.textContent="!~",ae=n(),ft=e("span"),ft.textContent=`NOT Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,ie=n(),D=e("li"),Pt=e("code"),Pt.textContent="?=",ot=n(),Rt=e("em"),Rt.textContent="Any/At least one of",S=n(),ut=e("span"),ut.textContent="Equal",re=n(),Q=e("li"),St=e("code"),St.textContent="?!=",Kt=n(),Ot=e("em"),Ot.textContent="Any/At least one of",ce=n(),N=e("span"),N.textContent="NOT equal",Vt=n(),z=e("li"),mt=e("code"),mt.textContent="?>",de=n(),I=e("em"),I.textContent="Any/At least one of",pe=n(),B=e("span"),B.textContent="Greater than",fe=n(),T=e("li"),Et=e("code"),Et.textContent="?>=",K=n(),ht=e("em"),ht.textContent="Any/At least one of",ue=n(),bt=e("span"),bt.textContent="Greater than or equal",me=n(),w=e("li"),Nt=e("code"),Nt.textContent="?<",at=n(),qt=e("em"),qt.textContent="Any/At least one of",he=n(),Mt=e("span"),Mt.textContent="Less than",Wt=n(),V=e("li"),gt=e("code"),gt.textContent="?<=",be=n(),Ht=e("em"),Ht.textContent="Any/At least one of",ge=n(),_t=e("span"),_t.textContent="Less than or equal",_e=n(),W=e("li"),xt=e("code"),xt.textContent="?~",Xt=n(),X=e("em"),X.textContent="Any/At least one of",Yt=n(),q=e("span"),q.textContent=`Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,wt=n(),y=e("li"),Dt=e("code"),Dt.textContent="?!~",xe=n(),Y=e("em"),Y.textContent="Any/At least one of",L=n(),It=e("span"),It.textContent=`NOT Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,G=n(),v=e("p"),v.innerHTML=`To group and combine several expressions you could use brackets
        <code>(...)</code>, <code>&amp;&amp;</code> (AND) and <code>||</code> (OR) tokens.`,a(b,"class","txt-danger"),a(Ct,"class","filter-op svelte-1w7s5nw"),a(O,"class","txt"),a(it,"class","filter-op svelte-1w7s5nw"),a(tt,"class","txt"),a(j,"class","filter-op svelte-1w7s5nw"),a(rt,"class","txt"),a(kt,"class","filter-op svelte-1w7s5nw"),a(ct,"class","txt"),a(E,"class","filter-op svelte-1w7s5nw"),a(yt,"class","txt"),a(vt,"class","filter-op svelte-1w7s5nw"),a(Ft,"class","txt"),a(Lt,"class","filter-op svelte-1w7s5nw"),a(At,"class","txt"),a(Tt,"class","filter-op svelte-1w7s5nw"),a(ft,"class","txt"),a(Pt,"class","filter-op svelte-1w7s5nw"),a(Rt,"class","txt-hint"),a(ut,"class","txt"),a(St,"class","filter-op svelte-1w7s5nw"),a(Ot,"class","txt-hint"),a(N,"class","txt"),a(mt,"class","filter-op svelte-1w7s5nw"),a(I,"class","txt-hint"),a(B,"class","txt"),a(Et,"class","filter-op svelte-1w7s5nw"),a(ht,"class","txt-hint"),a(bt,"class","txt"),a(Nt,"class","filter-op svelte-1w7s5nw"),a(qt,"class","txt-hint"),a(Mt,"class","txt"),a(gt,"class","filter-op svelte-1w7s5nw"),a(Ht,"class","txt-hint"),a(_t,"class","txt"),a(xt,"class","filter-op svelte-1w7s5nw"),a(X,"class","txt-hint"),a(q,"class","txt"),a(Dt,"class","filter-op svelte-1w7s5nw"),a(Y,"class","txt-hint"),a(It,"class","txt")},m(P,k){u(P,s,k),u(P,o,k),u(P,i,k),t(i,f),t(i,h),t(i,r),t(r,b),t(r,C),t(r,$),t(r,g),t(r,p),t(p,Z),t(Z,Ct),t(Z,Ut),t(Z,O),t(p,jt),t(p,H),t(H,it),t(H,R),t(H,tt),t(p,ne),t(p,U),t(U,j),t(U,se),t(U,rt),t(p,$t),t(p,et),t(et,kt),t(et,oe),t(et,ct),t(p,dt),t(p,lt),t(lt,E),t(lt,Jt),t(lt,yt),t(p,F),t(p,nt),t(nt,vt),t(nt,Qt),t(nt,Ft),t(p,J),t(p,st),t(st,Lt),t(st,zt),t(st,At),t(p,A),t(p,pt),t(pt,Tt),t(pt,ae),t(pt,ft),t(p,ie),t(p,D),t(D,Pt),t(D,ot),t(D,Rt),t(D,S),t(D,ut),t(p,re),t(p,Q),t(Q,St),t(Q,Kt),t(Q,Ot),t(Q,ce),t(Q,N),t(p,Vt),t(p,z),t(z,mt),t(z,de),t(z,I),t(z,pe),t(z,B),t(p,fe),t(p,T),t(T,Et),t(T,K),t(T,ht),t(T,ue),t(T,bt),t(p,me),t(p,w),t(w,Nt),t(w,at),t(w,qt),t(w,he),t(w,Mt),t(p,Wt),t(p,V),t(V,gt),t(V,be),t(V,Ht),t(V,ge),t(V,_t),t(p,_e),t(p,W),t(W,xt),t(W,Xt),t(W,X),t(W,Yt),t(W,q),t(p,wt),t(p,y),t(y,Dt),t(y,xe),t(y,Y),t(y,L),t(y,It),u(P,G,k),u(P,v,k)},d(P){P&&(m(s),m(o),m(i),m(G),m(v))}}}function fl(d){let s,o,i,f,h;function r(g,p){return g[0]?pl:dl}let b=r(d),C=b(d),$=d[0]&&Ue();return{c(){s=e("button"),C.c(),o=n(),$&&$.c(),i=el(),a(s,"class","btn btn-sm btn-secondary m-t-10")},m(g,p){u(g,s,p),C.m(s,null),u(g,o,p),$&&$.m(g,p),u(g,i,p),f||(h=tl(s,"click",d[1]),f=!0)},p(g,[p]){b!==(b=r(g))&&(C.d(1),C=b(g),C&&(C.c(),C.m(s,null))),g[0]?$||($=Ue(),$.c(),$.m(i.parentNode,i)):$&&($.d(1),$=null)},i:Be,o:Be,d(g){g&&(m(s),m(o),m(i)),C.d(),$&&$.d(g),f=!1,h()}}}function ul(d,s,o){let i=!1;function f(){o(0,i=!i)}return[i,f]}class ml extends Xe{constructor(s){super(),Ye(this,s,ul,fl,Ze,{})}}function je(d,s,o){const i=d.slice();return i[7]=s[o],i}function Je(d,s,o){const i=d.slice();return i[7]=s[o],i}function Qe(d,s,o){const i=d.slice();return i[12]=s[o],i[14]=o,i}function ze(d){let s;return{c(){s=e("p"),s.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",a(s,"class","txt-hint txt-sm txt-right")},m(o,i){u(o,s,i)},d(o){o&&m(s)}}}function Ke(d){let s,o=d[12]+"",i,f=d[14]<d[4].length-1?", ":"",h;return{c(){s=e("code"),i=_(o),h=_(f)},m(r,b){u(r,s,b),t(s,i),u(r,h,b)},p(r,b){b&16&&o!==(o=r[12]+"")&&Ce(i,o),b&16&&f!==(f=r[14]<r[4].length-1?", ":"")&&Ce(h,f)},d(r){r&&(m(s),m(h))}}}function Ve(d,s){let o,i,f;function h(){return s[6](s[7])}return{key:d,first:null,c(){o=e("button"),o.textContent=`${s[7].code} `,a(o,"type","button"),a(o,"class","tab-item"),ve(o,"active",s[2]===s[7].code),this.first=o},m(r,b){u(r,o,b),i||(f=tl(o,"click",h),i=!0)},p(r,b){s=r,b&36&&ve(o,"active",s[2]===s[7].code)},d(r){r&&m(o),i=!1,f()}}}function We(d,s){let o,i,f,h;return i=new ye({props:{content:s[7].body}}),{key:d,first:null,c(){o=e("div"),te(i.$$.fragment),f=n(),a(o,"class","tab-item"),ve(o,"active",s[2]===s[7].code),this.first=o},m(r,b){u(r,o,b),ee(i,o,null),t(o,f),h=!0},p(r,b){s=r,(!h||b&36)&&ve(o,"active",s[2]===s[7].code)},i(r){h||(Bt(i.$$.fragment,r),h=!0)},o(r){Gt(i.$$.fragment,r),h=!1},d(r){r&&m(o),le(i)}}}function hl(d){var Ae,Te,Pe,Re,Se,Oe;let s,o,i=d[0].name+"",f,h,r,b,C,$,g,p=d[0].name+"",Z,Ct,Ut,O,jt,H,it,R,tt,ne,U,j,se,rt,$t=d[0].name+"",et,kt,oe,ct,dt,lt,E,Jt,yt,F,nt,vt,Qt,Ft,J,st,Lt,zt,At,A,pt,Tt,ae,ft,ie,D,Pt,ot,Rt,S,ut,re,Q,St,Kt,Ot,ce,N,Vt,z,mt,de,I,pe,B,fe,T,Et,K,ht,ue,bt,me,w,Nt,at,qt,he,Mt,Wt,V,gt,be,Ht,ge,_t,_e,W,xt,Xt,X,Yt,q,wt,y=[],Dt=new Map,xe,Y,L=[],It=new Map,G;O=new cl({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${d[3]}');

        ...

        // fetch a paginated records list
        const resultList = await pb.collection('${(Ae=d[0])==null?void 0:Ae.name}').getList(1, 50, {
            filter: 'created >= "2022-01-01 00:00:00" && someField1 != someField2',
        });

        // you can also fetch all records at once via getFullList
        const records = await pb.collection('${(Te=d[0])==null?void 0:Te.name}').getFullList({
            sort: '-created',
        });

        // or fetch only the first record that matches the specified filter
        const record = await pb.collection('${(Pe=d[0])==null?void 0:Pe.name}').getFirstListItem('someField="test"', {
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
        final records = await pb.collection('${(Se=d[0])==null?void 0:Se.name}').getFullList(
          sort: '-created',
        );

        // or fetch only the first record that matches the specified filter
        final record = await pb.collection('${(Oe=d[0])==null?void 0:Oe.name}').getFirstListItem(
          'someField="test"',
          expand: 'relField1,relField2.subRelField',
        );
    `}});let v=d[1]&&ze();ot=new ye({props:{content:`
                        // DESC by created and ASC by id
                        ?sort=-created,id
                    `}});let P=Zt(d[4]),k=[];for(let l=0;l<P.length;l+=1)k[l]=Ke(Qe(d,P,l));B=new ye({props:{content:`
                        ?filter=(id='abc' && created>'2022-01-01')
                    `}}),T=new ml({}),at=new ye({props:{content:"?expand=relField1,relField2.subRelField"}});let $e=Zt(d[5]);const Fe=l=>l[7].code;for(let l=0;l<$e.length;l+=1){let c=Je(d,$e,l),x=Fe(c);Dt.set(x,y[l]=Ve(x,c))}let we=Zt(d[5]);const Le=l=>l[7].code;for(let l=0;l<we.length;l+=1){let c=je(d,we,l),x=Le(c);It.set(x,L[l]=We(x,c))}return{c(){s=e("h3"),o=_("List/Search ("),f=_(i),h=_(")"),r=n(),b=e("div"),C=e("p"),$=_("Fetch a paginated "),g=e("strong"),Z=_(p),Ct=_(" records list, supporting sorting and filtering."),Ut=n(),te(O.$$.fragment),jt=n(),H=e("h6"),H.textContent="API details",it=n(),R=e("div"),tt=e("strong"),tt.textContent="GET",ne=n(),U=e("div"),j=e("p"),se=_("/api/collections/"),rt=e("strong"),et=_($t),kt=_("/records"),oe=n(),v&&v.c(),ct=n(),dt=e("div"),dt.textContent="Query parameters",lt=n(),E=e("table"),Jt=e("thead"),Jt.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr>',yt=n(),F=e("tbody"),nt=e("tr"),nt.innerHTML='<td>page</td> <td><span class="label">Number</span></td> <td>The page (aka. offset) of the paginated list (default to 1).</td>',vt=n(),Qt=e("tr"),Qt.innerHTML='<td>perPage</td> <td><span class="label">Number</span></td> <td>Specify the max returned records per page (default to 30).</td>',Ft=n(),J=e("tr"),st=e("td"),st.textContent="sort",Lt=n(),zt=e("td"),zt.innerHTML='<span class="label">String</span>',At=n(),A=e("td"),pt=_("Specify the records order attribute(s). "),Tt=e("br"),ae=_(`
                Add `),ft=e("code"),ft.textContent="-",ie=_(" / "),D=e("code"),D.textContent="+",Pt=_(` (default) in front of the attribute for DESC / ASC order.
                Ex.:
                `),te(ot.$$.fragment),Rt=n(),S=e("p"),ut=e("strong"),ut.textContent="Supported record sort fields:",re=n(),Q=e("br"),St=n(),Kt=e("code"),Kt.textContent="@random",Ot=_(`,
                    `);for(let l=0;l<k.length;l+=1)k[l].c();ce=n(),N=e("tr"),Vt=e("td"),Vt.textContent="filter",z=n(),mt=e("td"),mt.innerHTML='<span class="label">String</span>',de=n(),I=e("td"),pe=_(`Filter the returned records. Ex.:
                `),te(B.$$.fragment),fe=n(),te(T.$$.fragment),Et=n(),K=e("tr"),ht=e("td"),ht.textContent="expand",ue=n(),bt=e("td"),bt.innerHTML='<span class="label">String</span>',me=n(),w=e("td"),Nt=_(`Auto expand record relations. Ex.:
                `),te(at.$$.fragment),qt=_(`
                Supports up to 6-levels depth nested relations expansion. `),he=e("br"),Mt=_(`
                The expanded relations will be appended to each individual record under the
                `),Wt=e("code"),Wt.textContent="expand",V=_(" property (eg. "),gt=e("code"),gt.textContent='"expand": {"relField1": {...}, ...}',be=_(`).
                `),Ht=e("br"),ge=_(`
                Only the relations to which the request user has permissions to `),_t=e("strong"),_t.textContent="view",_e=_(" will be expanded."),W=n(),xt=e("tr"),xt.innerHTML=`<td id="query-page">fields</td> <td><span class="label">String</span></td> <td>Comma separated string of the fields to return in the JSON response
                <em>(by default returns all fields)</em>.</td>`,Xt=n(),X=e("div"),X.textContent="Responses",Yt=n(),q=e("div"),wt=e("div");for(let l=0;l<y.length;l+=1)y[l].c();xe=n(),Y=e("div");for(let l=0;l<L.length;l+=1)L[l].c();a(s,"class","m-b-sm"),a(b,"class","content txt-lg m-b-sm"),a(H,"class","m-b-xs"),a(tt,"class","label label-primary"),a(U,"class","content"),a(R,"class","alert alert-info"),a(dt,"class","section-title"),a(E,"class","table-compact table-border m-b-base"),a(X,"class","section-title"),a(wt,"class","tabs-header compact left"),a(Y,"class","tabs-content"),a(q,"class","tabs")},m(l,c){u(l,s,c),t(s,o),t(s,f),t(s,h),u(l,r,c),u(l,b,c),t(b,C),t(C,$),t(C,g),t(g,Z),t(C,Ct),u(l,Ut,c),ee(O,l,c),u(l,jt,c),u(l,H,c),u(l,it,c),u(l,R,c),t(R,tt),t(R,ne),t(R,U),t(U,j),t(j,se),t(j,rt),t(rt,et),t(j,kt),t(R,oe),v&&v.m(R,null),u(l,ct,c),u(l,dt,c),u(l,lt,c),u(l,E,c),t(E,Jt),t(E,yt),t(E,F),t(F,nt),t(F,vt),t(F,Qt),t(F,Ft),t(F,J),t(J,st),t(J,Lt),t(J,zt),t(J,At),t(J,A),t(A,pt),t(A,Tt),t(A,ae),t(A,ft),t(A,ie),t(A,D),t(A,Pt),ee(ot,A,null),t(A,Rt),t(A,S),t(S,ut),t(S,re),t(S,Q),t(S,St),t(S,Kt),t(S,Ot);for(let x=0;x<k.length;x+=1)k[x]&&k[x].m(S,null);t(F,ce),t(F,N),t(N,Vt),t(N,z),t(N,mt),t(N,de),t(N,I),t(I,pe),ee(B,I,null),t(I,fe),ee(T,I,null),t(F,Et),t(F,K),t(K,ht),t(K,ue),t(K,bt),t(K,me),t(K,w),t(w,Nt),ee(at,w,null),t(w,qt),t(w,he),t(w,Mt),t(w,Wt),t(w,V),t(w,gt),t(w,be),t(w,Ht),t(w,ge),t(w,_t),t(w,_e),t(F,W),t(F,xt),u(l,Xt,c),u(l,X,c),u(l,Yt,c),u(l,q,c),t(q,wt);for(let x=0;x<y.length;x+=1)y[x]&&y[x].m(wt,null);t(q,xe),t(q,Y);for(let x=0;x<L.length;x+=1)L[x]&&L[x].m(Y,null);G=!0},p(l,[c]){var Ee,Ne,qe,Me,He,De;(!G||c&1)&&i!==(i=l[0].name+"")&&Ce(f,i),(!G||c&1)&&p!==(p=l[0].name+"")&&Ce(Z,p);const x={};if(c&9&&(x.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${l[3]}');

        ...

        // fetch a paginated records list
        const resultList = await pb.collection('${(Ee=l[0])==null?void 0:Ee.name}').getList(1, 50, {
            filter: 'created >= "2022-01-01 00:00:00" && someField1 != someField2',
        });

        // you can also fetch all records at once via getFullList
        const records = await pb.collection('${(Ne=l[0])==null?void 0:Ne.name}').getFullList({
            sort: '-created',
        });

        // or fetch only the first record that matches the specified filter
        const record = await pb.collection('${(qe=l[0])==null?void 0:qe.name}').getFirstListItem('someField="test"', {
            expand: 'relField1,relField2.subRelField',
        });
    `),c&9&&(x.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${l[3]}');

        ...

        // fetch a paginated records list
        final resultList = await pb.collection('${(Me=l[0])==null?void 0:Me.name}').getList(
          page: 1,
          perPage: 50,
          filter: 'created >= "2022-01-01 00:00:00" && someField1 != someField2',
        );

        // you can also fetch all records at once via getFullList
        final records = await pb.collection('${(He=l[0])==null?void 0:He.name}').getFullList(
          sort: '-created',
        );

        // or fetch only the first record that matches the specified filter
        final record = await pb.collection('${(De=l[0])==null?void 0:De.name}').getFirstListItem(
          'someField="test"',
          expand: 'relField1,relField2.subRelField',
        );
    `),O.$set(x),(!G||c&1)&&$t!==($t=l[0].name+"")&&Ce(et,$t),l[1]?v||(v=ze(),v.c(),v.m(R,null)):v&&(v.d(1),v=null),c&16){P=Zt(l[4]);let M;for(M=0;M<P.length;M+=1){const Ie=Qe(l,P,M);k[M]?k[M].p(Ie,c):(k[M]=Ke(Ie),k[M].c(),k[M].m(S,null))}for(;M<k.length;M+=1)k[M].d(1);k.length=P.length}c&36&&($e=Zt(l[5]),y=Ge(y,c,Fe,1,l,$e,Dt,wt,ll,Ve,null,Je)),c&36&&(we=Zt(l[5]),nl(),L=Ge(L,c,Le,1,l,we,It,Y,sl,We,null,je),ol())},i(l){if(!G){Bt(O.$$.fragment,l),Bt(ot.$$.fragment,l),Bt(B.$$.fragment,l),Bt(T.$$.fragment,l),Bt(at.$$.fragment,l);for(let c=0;c<we.length;c+=1)Bt(L[c]);G=!0}},o(l){Gt(O.$$.fragment,l),Gt(ot.$$.fragment,l),Gt(B.$$.fragment,l),Gt(T.$$.fragment,l),Gt(at.$$.fragment,l);for(let c=0;c<L.length;c+=1)Gt(L[c]);G=!1},d(l){l&&(m(s),m(r),m(b),m(Ut),m(jt),m(H),m(it),m(R),m(ct),m(dt),m(lt),m(E),m(Xt),m(X),m(Yt),m(q)),le(O,l),v&&v.d(),le(ot),al(k,l),le(B),le(T),le(at);for(let c=0;c<y.length;c+=1)y[c].d();for(let c=0;c<L.length;c+=1)L[c].d()}}}function bl(d,s,o){let i,f,h,{collection:r=new il}=s,b=200,C=[];const $=g=>o(2,b=g.code);return d.$$set=g=>{"collection"in g&&o(0,r=g.collection)},d.$$.update=()=>{d.$$.dirty&1&&o(4,i=ke.getAllCollectionIdentifiers(r)),d.$$.dirty&1&&o(1,f=(r==null?void 0:r.listRule)===null),d.$$.dirty&3&&r!=null&&r.id&&(C.push({code:200,body:JSON.stringify({page:1,perPage:30,totalPages:1,totalItems:2,items:[ke.dummyCollectionRecord(r),ke.dummyCollectionRecord(r)]},null,2)}),C.push({code:400,body:`
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
                `}))},o(3,h=ke.getApiExampleUrl(rl.baseUrl)),[r,f,b,h,i,C,$]}class xl extends Xe{constructor(s){super(),Ye(this,s,bl,hl,Ze,{collection:0})}}export{xl as default};
