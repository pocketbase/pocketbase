import{S as Ce,i as Re,s as Pe,e as c,w as D,b as k,c as $e,f as m,g as d,h as n,m as we,x,O as _e,P as Ee,k as Oe,Q as Te,n as Be,t as ee,a as te,o as f,d as ge,R as Ie,C as Ae,p as Me,r as N,u as Se,N as qe}from"./index.89a3f554.js";import{S as He}from"./SdkTabs.0a6ad1c9.js";function ke(o,l,s){const a=o.slice();return a[6]=l[s],a}function he(o,l,s){const a=o.slice();return a[6]=l[s],a}function ve(o){let l;return{c(){l=c("p"),l.innerHTML="Requires admin <code>Authorization:TOKEN</code> header",m(l,"class","txt-hint txt-sm txt-right")},m(s,a){d(s,l,a)},d(s){s&&f(l)}}}function ye(o,l){let s,a=l[6].code+"",h,i,r,u;function $(){return l[5](l[6])}return{key:o,first:null,c(){s=c("button"),h=D(a),i=k(),m(s,"class","tab-item"),N(s,"active",l[2]===l[6].code),this.first=s},m(b,g){d(b,s,g),n(s,h),n(s,i),r||(u=Se(s,"click",$),r=!0)},p(b,g){l=b,g&20&&N(s,"active",l[2]===l[6].code)},d(b){b&&f(s),r=!1,u()}}}function De(o,l){let s,a,h,i;return a=new qe({props:{content:l[6].body}}),{key:o,first:null,c(){s=c("div"),$e(a.$$.fragment),h=k(),m(s,"class","tab-item"),N(s,"active",l[2]===l[6].code),this.first=s},m(r,u){d(r,s,u),we(a,s,null),n(s,h),i=!0},p(r,u){l=r,(!i||u&20)&&N(s,"active",l[2]===l[6].code)},i(r){i||(ee(a.$$.fragment,r),i=!0)},o(r){te(a.$$.fragment,r),i=!1},d(r){r&&f(s),ge(a)}}}function Le(o){var ue,pe;let l,s,a=o[0].name+"",h,i,r,u,$,b,g,q=o[0].name+"",z,le,F,C,K,O,Q,y,H,se,L,E,oe,G,U=o[0].name+"",J,ae,V,ne,W,T,X,B,Y,I,Z,R,A,w=[],ie=new Map,re,M,v=[],ce=new Map,P;C=new He({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${o[3]}');

        ...

        await pb.collection('${(ue=o[0])==null?void 0:ue.name}').delete('RECORD_ID');
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${o[3]}');

        ...

        await pb.collection('${(pe=o[0])==null?void 0:pe.name}').delete('RECORD_ID');
    `}});let _=o[1]&&ve(),j=o[4];const de=e=>e[6].code;for(let e=0;e<j.length;e+=1){let t=he(o,j,e),p=de(t);ie.set(p,w[e]=ye(p,t))}let S=o[4];const fe=e=>e[6].code;for(let e=0;e<S.length;e+=1){let t=ke(o,S,e),p=fe(t);ce.set(p,v[e]=De(p,t))}return{c(){l=c("h3"),s=D("Delete ("),h=D(a),i=D(")"),r=k(),u=c("div"),$=c("p"),b=D("Delete a single "),g=c("strong"),z=D(q),le=D(" record."),F=k(),$e(C.$$.fragment),K=k(),O=c("h6"),O.textContent="API details",Q=k(),y=c("div"),H=c("strong"),H.textContent="DELETE",se=k(),L=c("div"),E=c("p"),oe=D("/api/collections/"),G=c("strong"),J=D(U),ae=D("/records/"),V=c("strong"),V.textContent=":id",ne=k(),_&&_.c(),W=k(),T=c("div"),T.textContent="Path parameters",X=k(),B=c("table"),B.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr></thead> 
    <tbody><tr><td>id</td> 
            <td><span class="label">String</span></td> 
            <td>ID of the record to delete.</td></tr></tbody>`,Y=k(),I=c("div"),I.textContent="Responses",Z=k(),R=c("div"),A=c("div");for(let e=0;e<w.length;e+=1)w[e].c();re=k(),M=c("div");for(let e=0;e<v.length;e+=1)v[e].c();m(l,"class","m-b-sm"),m(u,"class","content txt-lg m-b-sm"),m(O,"class","m-b-xs"),m(H,"class","label label-primary"),m(L,"class","content"),m(y,"class","alert alert-danger"),m(T,"class","section-title"),m(B,"class","table-compact table-border m-b-base"),m(I,"class","section-title"),m(A,"class","tabs-header compact left"),m(M,"class","tabs-content"),m(R,"class","tabs")},m(e,t){d(e,l,t),n(l,s),n(l,h),n(l,i),d(e,r,t),d(e,u,t),n(u,$),n($,b),n($,g),n(g,z),n($,le),d(e,F,t),we(C,e,t),d(e,K,t),d(e,O,t),d(e,Q,t),d(e,y,t),n(y,H),n(y,se),n(y,L),n(L,E),n(E,oe),n(E,G),n(G,J),n(E,ae),n(E,V),n(y,ne),_&&_.m(y,null),d(e,W,t),d(e,T,t),d(e,X,t),d(e,B,t),d(e,Y,t),d(e,I,t),d(e,Z,t),d(e,R,t),n(R,A);for(let p=0;p<w.length;p+=1)w[p].m(A,null);n(R,re),n(R,M);for(let p=0;p<v.length;p+=1)v[p].m(M,null);P=!0},p(e,[t]){var me,be;(!P||t&1)&&a!==(a=e[0].name+"")&&x(h,a),(!P||t&1)&&q!==(q=e[0].name+"")&&x(z,q);const p={};t&9&&(p.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(me=e[0])==null?void 0:me.name}').delete('RECORD_ID');
    `),t&9&&(p.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(be=e[0])==null?void 0:be.name}').delete('RECORD_ID');
    `),C.$set(p),(!P||t&1)&&U!==(U=e[0].name+"")&&x(J,U),e[1]?_||(_=ve(),_.c(),_.m(y,null)):_&&(_.d(1),_=null),t&20&&(j=e[4],w=_e(w,t,de,1,e,j,ie,A,Ee,ye,null,he)),t&20&&(S=e[4],Oe(),v=_e(v,t,fe,1,e,S,ce,M,Te,De,null,ke),Be())},i(e){if(!P){ee(C.$$.fragment,e);for(let t=0;t<S.length;t+=1)ee(v[t]);P=!0}},o(e){te(C.$$.fragment,e);for(let t=0;t<v.length;t+=1)te(v[t]);P=!1},d(e){e&&f(l),e&&f(r),e&&f(u),e&&f(F),ge(C,e),e&&f(K),e&&f(O),e&&f(Q),e&&f(y),_&&_.d(),e&&f(W),e&&f(T),e&&f(X),e&&f(B),e&&f(Y),e&&f(I),e&&f(Z),e&&f(R);for(let t=0;t<w.length;t+=1)w[t].d();for(let t=0;t<v.length;t+=1)v[t].d()}}}function Ue(o,l,s){let a,h,{collection:i=new Ie}=l,r=204,u=[];const $=b=>s(2,r=b.code);return o.$$set=b=>{"collection"in b&&s(0,i=b.collection)},o.$$.update=()=>{o.$$.dirty&1&&s(1,a=(i==null?void 0:i.deleteRule)===null),o.$$.dirty&3&&i!=null&&i.id&&(u.push({code:204,body:`
                null
            `}),u.push({code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to delete record. Make sure that the record is not part of a required relation reference.",
                  "data": {}
                }
            `}),a&&u.push({code:403,body:`
                    {
                      "code": 403,
                      "message": "Only admins can access this action.",
                      "data": {}
                    }
                `}),u.push({code:404,body:`
                {
                  "code": 404,
                  "message": "The requested resource wasn't found.",
                  "data": {}
                }
            `}))},s(3,h=Ae.getApiExampleUrl(Me.baseUrl)),[i,a,r,h,u,$]}class ze extends Ce{constructor(l){super(),Re(this,l,Ue,Le,Pe,{collection:0})}}export{ze as default};
