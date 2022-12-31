import{S as Et,i as Nt,s as Ht,e as l,b as a,E as qt,f as d,g as p,u as Mt,y as xt,o as u,w as k,h as e,N as Ae,c as ge,m as ye,x as Ue,O as Lt,P as Dt,k as It,Q as Bt,n as zt,t as ce,a as de,d as ve,R as Gt,C as je,p as Ut,r as Ee}from"./index.89a3f554.js";import{S as jt}from"./SdkTabs.0a6ad1c9.js";function Qt(r){let s,n,i;return{c(){s=l("span"),s.textContent="Show details",n=a(),i=l("i"),d(s,"class","txt"),d(i,"class","ri-arrow-down-s-line")},m(c,f){p(c,s,f),p(c,n,f),p(c,i,f)},d(c){c&&u(s),c&&u(n),c&&u(i)}}}function Jt(r){let s,n,i;return{c(){s=l("span"),s.textContent="Hide details",n=a(),i=l("i"),d(s,"class","txt"),d(i,"class","ri-arrow-up-s-line")},m(c,f){p(c,s,f),p(c,n,f),p(c,i,f)},d(c){c&&u(s),c&&u(n),c&&u(i)}}}function Tt(r){let s,n,i,c,f,m,_,w,b,$,h,H,W,fe,T,pe,O,G,C,M,Fe,A,E,Ce,U,X,q,Y,xe,j,Q,D,P,ue,Z,v,I,ee,me,te,N,B,le,be,se,x,J,ne,Le,K,he,V;return{c(){s=l("p"),s.innerHTML=`The syntax basically follows the format
        <code><span class="txt-success">OPERAND</span> 
            <span class="txt-danger">OPERATOR</span> 
            <span class="txt-success">OPERAND</span></code>, where:`,n=a(),i=l("ul"),c=l("li"),c.innerHTML=`<code class="txt-success">OPERAND</code> - could be any of the above field literal, string (single
            or double quoted), number, null, true, false`,f=a(),m=l("li"),_=l("code"),_.textContent="OPERATOR",w=k(` - is one of:
            `),b=l("br"),$=a(),h=l("ul"),H=l("li"),W=l("code"),W.textContent="=",fe=a(),T=l("span"),T.textContent="Equal",pe=a(),O=l("li"),G=l("code"),G.textContent="!=",C=a(),M=l("span"),M.textContent="NOT equal",Fe=a(),A=l("li"),E=l("code"),E.textContent=">",Ce=a(),U=l("span"),U.textContent="Greater than",X=a(),q=l("li"),Y=l("code"),Y.textContent=">=",xe=a(),j=l("span"),j.textContent="Greater than or equal",Q=a(),D=l("li"),P=l("code"),P.textContent="<",ue=a(),Z=l("span"),Z.textContent="Less than or equal",v=a(),I=l("li"),ee=l("code"),ee.textContent="<=",me=a(),te=l("span"),te.textContent="Less than or equal",N=a(),B=l("li"),le=l("code"),le.textContent="~",be=a(),se=l("span"),se.textContent=`Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,x=a(),J=l("li"),ne=l("code"),ne.textContent="!~",Le=a(),K=l("span"),K.textContent=`NOT Like/Contains (if not specified auto wraps the right string OPERAND in a "%" for
                        wildcard match)`,he=a(),V=l("p"),V.innerHTML=`To group and combine several expressions you could use brackets
        <code>(...)</code>, <code>&amp;&amp;</code> (AND) and <code>||</code> (OR) tokens.`,d(_,"class","txt-danger"),d(W,"class","filter-op svelte-1w7s5nw"),d(T,"class","txt-hint"),d(G,"class","filter-op svelte-1w7s5nw"),d(M,"class","txt-hint"),d(E,"class","filter-op svelte-1w7s5nw"),d(U,"class","txt-hint"),d(Y,"class","filter-op svelte-1w7s5nw"),d(j,"class","txt-hint"),d(P,"class","filter-op svelte-1w7s5nw"),d(Z,"class","txt-hint"),d(ee,"class","filter-op svelte-1w7s5nw"),d(te,"class","txt-hint"),d(le,"class","filter-op svelte-1w7s5nw"),d(se,"class","txt-hint"),d(ne,"class","filter-op svelte-1w7s5nw"),d(K,"class","txt-hint")},m(F,R){p(F,s,R),p(F,n,R),p(F,i,R),e(i,c),e(i,f),e(i,m),e(m,_),e(m,w),e(m,b),e(m,$),e(m,h),e(h,H),e(H,W),e(H,fe),e(H,T),e(h,pe),e(h,O),e(O,G),e(O,C),e(O,M),e(h,Fe),e(h,A),e(A,E),e(A,Ce),e(A,U),e(h,X),e(h,q),e(q,Y),e(q,xe),e(q,j),e(h,Q),e(h,D),e(D,P),e(D,ue),e(D,Z),e(h,v),e(h,I),e(I,ee),e(I,me),e(I,te),e(h,N),e(h,B),e(B,le),e(B,be),e(B,se),e(h,x),e(h,J),e(J,ne),e(J,Le),e(J,K),p(F,he,R),p(F,V,R)},d(F){F&&u(s),F&&u(n),F&&u(i),F&&u(he),F&&u(V)}}}function Kt(r){let s,n,i,c,f;function m($,h){return $[0]?Jt:Qt}let _=m(r),w=_(r),b=r[0]&&Tt();return{c(){s=l("button"),w.c(),n=a(),b&&b.c(),i=qt(),d(s,"class","btn btn-sm btn-secondary m-t-5")},m($,h){p($,s,h),w.m(s,null),p($,n,h),b&&b.m($,h),p($,i,h),c||(f=Mt(s,"click",r[1]),c=!0)},p($,[h]){_!==(_=m($))&&(w.d(1),w=_($),w&&(w.c(),w.m(s,null))),$[0]?b||(b=Tt(),b.c(),b.m(i.parentNode,i)):b&&(b.d(1),b=null)},i:xt,o:xt,d($){$&&u(s),w.d(),$&&u(n),b&&b.d($),$&&u(i),c=!1,f()}}}function Vt(r,s,n){let i=!1;function c(){n(0,i=!i)}return[i,c]}class Wt extends Et{constructor(s){super(),Nt(this,s,Vt,Kt,Ht,{})}}function Pt(r,s,n){const i=r.slice();return i[6]=s[n],i}function Rt(r,s,n){const i=r.slice();return i[6]=s[n],i}function St(r){let s;return{c(){s=l("p"),s.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",d(s,"class","txt-hint txt-sm txt-right")},m(n,i){p(n,s,i)},d(n){n&&u(s)}}}function Ot(r,s){let n,i=s[6].code+"",c,f,m,_;function w(){return s[5](s[6])}return{key:r,first:null,c(){n=l("div"),c=k(i),f=a(),d(n,"class","tab-item"),Ee(n,"active",s[2]===s[6].code),this.first=n},m(b,$){p(b,n,$),e(n,c),e(n,f),m||(_=Mt(n,"click",w),m=!0)},p(b,$){s=b,$&20&&Ee(n,"active",s[2]===s[6].code)},d(b){b&&u(n),m=!1,_()}}}function At(r,s){let n,i,c,f;return i=new Ae({props:{content:s[6].body}}),{key:r,first:null,c(){n=l("div"),ge(i.$$.fragment),c=a(),d(n,"class","tab-item"),Ee(n,"active",s[2]===s[6].code),this.first=n},m(m,_){p(m,n,_),ye(i,n,null),e(n,c),f=!0},p(m,_){s=m,(!f||_&20)&&Ee(n,"active",s[2]===s[6].code)},i(m){f||(ce(i.$$.fragment,m),f=!0)},o(m){de(i.$$.fragment,m),f=!1},d(m){m&&u(n),ve(i)}}}function Xt(r){var mt,bt,ht,_t,$t,kt;let s,n,i=r[0].name+"",c,f,m,_,w,b,$,h=r[0].name+"",H,W,fe,T,pe,O,G,C,M,Fe,A,E,Ce,U,X=r[0].name+"",q,Y,xe,j,Q,D,P,ue,Z,v,I,ee,me,te,N,B,le,be,se,x,J,ne,Le,K,he,V,F,R,Qe,ie,Ne,Je,He,Ke,_e,Ve,$e,We,ke,Xe,oe,Me,Ye,qe,Ze,y,et,we,tt,lt,st,De,nt,Ie,it,ot,at,Be,rt,ze,Te,Ge,ae,Pe,z=[],ct=new Map,dt,Re,S=[],ft=new Map,re;T=new jt({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${r[3]}');

        ...

        // fetch a paginated records list
        const resultList = await pb.collection('${(mt=r[0])==null?void 0:mt.name}').getList(1, 50, {
            filter: 'created >= "2022-01-01 00:00:00" && someFiled1 != someField2',
        });

        // you can also fetch all records at once via getFullList
        const records = await pb.collection('${(bt=r[0])==null?void 0:bt.name}').getFullList(200 /* batch size */, {
            sort: '-created',
        });

        // or fetch only the first record that matches the specified filter
        const record = await pb.collection('${(ht=r[0])==null?void 0:ht.name}').getFirstListItem('someField="test"', {
            expand: 'relField1,relField2.subRelField',
        });
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${r[3]}');

        ...

        // fetch a paginated records list
        final resultList = await pb.collection('${(_t=r[0])==null?void 0:_t.name}').getList(
          page: 1,
          perPage: 50,
          filter: 'created >= "2022-01-01 00:00:00" && someFiled1 != someField2',
        );

        // you can also fetch all records at once via getFullList
        final records = await pb.collection('${($t=r[0])==null?void 0:$t.name}').getFullList(
          batch: 200,
          sort: '-created',
        );

        // or fetch only the first record that matches the specified filter
        final record = await pb.collection('${(kt=r[0])==null?void 0:kt.name}').getFirstListItem(
          'someField="test"',
          expand: 'relField1,relField2.subRelField',
        );
    `}});let L=r[1]&&St();R=new Ae({props:{content:`
                        // DESC by created and ASC by id
                        ?sort=-created,id
                    `}}),$e=new Ae({props:{content:`
                        ?filter=(id='abc' && created>'2022-01-01')
                    `}}),ke=new Wt({}),we=new Ae({props:{content:"?expand=relField1,relField2.subRelField"}});let Oe=r[4];const pt=t=>t[6].code;for(let t=0;t<Oe.length;t+=1){let o=Rt(r,Oe,t),g=pt(o);ct.set(g,z[t]=Ot(g,o))}let Se=r[4];const ut=t=>t[6].code;for(let t=0;t<Se.length;t+=1){let o=Pt(r,Se,t),g=ut(o);ft.set(g,S[t]=At(g,o))}return{c(){s=l("h3"),n=k("List/Search ("),c=k(i),f=k(")"),m=a(),_=l("div"),w=l("p"),b=k("Fetch a paginated "),$=l("strong"),H=k(h),W=k(" records list, supporting sorting and filtering."),fe=a(),ge(T.$$.fragment),pe=a(),O=l("h6"),O.textContent="API details",G=a(),C=l("div"),M=l("strong"),M.textContent="GET",Fe=a(),A=l("div"),E=l("p"),Ce=k("/api/collections/"),U=l("strong"),q=k(X),Y=k("/records"),xe=a(),L&&L.c(),j=a(),Q=l("div"),Q.textContent="Query parameters",D=a(),P=l("table"),ue=l("thead"),ue.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,Z=a(),v=l("tbody"),I=l("tr"),I.innerHTML=`<td>page</td> 
            <td><span class="label">Number</span></td> 
            <td>The page (aka. offset) of the paginated list (default to 1).</td>`,ee=a(),me=l("tr"),me.innerHTML=`<td>perPage</td> 
            <td><span class="label">Number</span></td> 
            <td>Specify the max returned records per page (default to 30).</td>`,te=a(),N=l("tr"),B=l("td"),B.textContent="sort",le=a(),be=l("td"),be.innerHTML='<span class="label">String</span>',se=a(),x=l("td"),J=k("Specify the records order attribute(s). "),ne=l("br"),Le=k(`
                Add `),K=l("code"),K.textContent="-",he=k(" / "),V=l("code"),V.textContent="+",F=k(` (default) in front of the attribute for DESC / ASC order.
                Ex.:
                `),ge(R.$$.fragment),Qe=a(),ie=l("tr"),Ne=l("td"),Ne.textContent="filter",Je=a(),He=l("td"),He.innerHTML='<span class="label">String</span>',Ke=a(),_e=l("td"),Ve=k(`Filter the returned records. Ex.:
                `),ge($e.$$.fragment),We=a(),ge(ke.$$.fragment),Xe=a(),oe=l("tr"),Me=l("td"),Me.textContent="expand",Ye=a(),qe=l("td"),qe.innerHTML='<span class="label">String</span>',Ze=a(),y=l("td"),et=k(`Auto expand record relations. Ex.:
                `),ge(we.$$.fragment),tt=k(`
                Supports up to 6-levels depth nested relations expansion. `),lt=l("br"),st=k(`
                The expanded relations will be appended to each individual record under the
                `),De=l("code"),De.textContent="expand",nt=k(" property (eg. "),Ie=l("code"),Ie.textContent='"expand": {"relField1": {...}, ...}',it=k(`).
                `),ot=l("br"),at=k(`
                Only the relations to which the request user has permissions to `),Be=l("strong"),Be.textContent="view",rt=k(" will be expanded."),ze=a(),Te=l("div"),Te.textContent="Responses",Ge=a(),ae=l("div"),Pe=l("div");for(let t=0;t<z.length;t+=1)z[t].c();dt=a(),Re=l("div");for(let t=0;t<S.length;t+=1)S[t].c();d(s,"class","m-b-sm"),d(_,"class","content txt-lg m-b-sm"),d(O,"class","m-b-xs"),d(M,"class","label label-primary"),d(A,"class","content"),d(C,"class","alert alert-info"),d(Q,"class","section-title"),d(P,"class","table-compact table-border m-b-base"),d(Te,"class","section-title"),d(Pe,"class","tabs-header compact left"),d(Re,"class","tabs-content"),d(ae,"class","tabs")},m(t,o){p(t,s,o),e(s,n),e(s,c),e(s,f),p(t,m,o),p(t,_,o),e(_,w),e(w,b),e(w,$),e($,H),e(w,W),p(t,fe,o),ye(T,t,o),p(t,pe,o),p(t,O,o),p(t,G,o),p(t,C,o),e(C,M),e(C,Fe),e(C,A),e(A,E),e(E,Ce),e(E,U),e(U,q),e(E,Y),e(C,xe),L&&L.m(C,null),p(t,j,o),p(t,Q,o),p(t,D,o),p(t,P,o),e(P,ue),e(P,Z),e(P,v),e(v,I),e(v,ee),e(v,me),e(v,te),e(v,N),e(N,B),e(N,le),e(N,be),e(N,se),e(N,x),e(x,J),e(x,ne),e(x,Le),e(x,K),e(x,he),e(x,V),e(x,F),ye(R,x,null),e(v,Qe),e(v,ie),e(ie,Ne),e(ie,Je),e(ie,He),e(ie,Ke),e(ie,_e),e(_e,Ve),ye($e,_e,null),e(_e,We),ye(ke,_e,null),e(v,Xe),e(v,oe),e(oe,Me),e(oe,Ye),e(oe,qe),e(oe,Ze),e(oe,y),e(y,et),ye(we,y,null),e(y,tt),e(y,lt),e(y,st),e(y,De),e(y,nt),e(y,Ie),e(y,it),e(y,ot),e(y,at),e(y,Be),e(y,rt),p(t,ze,o),p(t,Te,o),p(t,Ge,o),p(t,ae,o),e(ae,Pe);for(let g=0;g<z.length;g+=1)z[g].m(Pe,null);e(ae,dt),e(ae,Re);for(let g=0;g<S.length;g+=1)S[g].m(Re,null);re=!0},p(t,[o]){var wt,gt,yt,vt,Ft,Ct;(!re||o&1)&&i!==(i=t[0].name+"")&&Ue(c,i),(!re||o&1)&&h!==(h=t[0].name+"")&&Ue(H,h);const g={};o&9&&(g.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${t[3]}');

        ...

        // fetch a paginated records list
        const resultList = await pb.collection('${(wt=t[0])==null?void 0:wt.name}').getList(1, 50, {
            filter: 'created >= "2022-01-01 00:00:00" && someFiled1 != someField2',
        });

        // you can also fetch all records at once via getFullList
        const records = await pb.collection('${(gt=t[0])==null?void 0:gt.name}').getFullList(200 /* batch size */, {
            sort: '-created',
        });

        // or fetch only the first record that matches the specified filter
        const record = await pb.collection('${(yt=t[0])==null?void 0:yt.name}').getFirstListItem('someField="test"', {
            expand: 'relField1,relField2.subRelField',
        });
    `),o&9&&(g.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${t[3]}');

        ...

        // fetch a paginated records list
        final resultList = await pb.collection('${(vt=t[0])==null?void 0:vt.name}').getList(
          page: 1,
          perPage: 50,
          filter: 'created >= "2022-01-01 00:00:00" && someFiled1 != someField2',
        );

        // you can also fetch all records at once via getFullList
        final records = await pb.collection('${(Ft=t[0])==null?void 0:Ft.name}').getFullList(
          batch: 200,
          sort: '-created',
        );

        // or fetch only the first record that matches the specified filter
        final record = await pb.collection('${(Ct=t[0])==null?void 0:Ct.name}').getFirstListItem(
          'someField="test"',
          expand: 'relField1,relField2.subRelField',
        );
    `),T.$set(g),(!re||o&1)&&X!==(X=t[0].name+"")&&Ue(q,X),t[1]?L||(L=St(),L.c(),L.m(C,null)):L&&(L.d(1),L=null),o&20&&(Oe=t[4],z=Lt(z,o,pt,1,t,Oe,ct,Pe,Dt,Ot,null,Rt)),o&20&&(Se=t[4],It(),S=Lt(S,o,ut,1,t,Se,ft,Re,Bt,At,null,Pt),zt())},i(t){if(!re){ce(T.$$.fragment,t),ce(R.$$.fragment,t),ce($e.$$.fragment,t),ce(ke.$$.fragment,t),ce(we.$$.fragment,t);for(let o=0;o<Se.length;o+=1)ce(S[o]);re=!0}},o(t){de(T.$$.fragment,t),de(R.$$.fragment,t),de($e.$$.fragment,t),de(ke.$$.fragment,t),de(we.$$.fragment,t);for(let o=0;o<S.length;o+=1)de(S[o]);re=!1},d(t){t&&u(s),t&&u(m),t&&u(_),t&&u(fe),ve(T,t),t&&u(pe),t&&u(O),t&&u(G),t&&u(C),L&&L.d(),t&&u(j),t&&u(Q),t&&u(D),t&&u(P),ve(R),ve($e),ve(ke),ve(we),t&&u(ze),t&&u(Te),t&&u(Ge),t&&u(ae);for(let o=0;o<z.length;o+=1)z[o].d();for(let o=0;o<S.length;o+=1)S[o].d()}}}function Yt(r,s,n){let i,c,{collection:f=new Gt}=s,m=200,_=[];const w=b=>n(2,m=b.code);return r.$$set=b=>{"collection"in b&&n(0,f=b.collection)},r.$$.update=()=>{r.$$.dirty&1&&n(1,i=(f==null?void 0:f.listRule)===null),r.$$.dirty&3&&f!=null&&f.id&&(_.push({code:200,body:JSON.stringify({page:1,perPage:30,totalPages:1,totalItems:2,items:[je.dummyCollectionRecord(f),je.dummyCollectionRecord(f)]},null,2)}),_.push({code:400,body:`
                {
                  "code": 400,
                  "message": "Something went wrong while processing your request. Invalid filter.",
                  "data": {}
                }
            `}),i&&_.push({code:403,body:`
                    {
                      "code": 403,
                      "message": "Only admins can access this action.",
                      "data": {}
                    }
                `}))},n(3,c=je.getApiExampleUrl(Ut.baseUrl)),[f,i,m,c,_,w]}class tl extends Et{constructor(s){super(),Nt(this,s,Yt,Xt,Ht,{collection:0})}}export{tl as default};
